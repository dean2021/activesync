package ashttp

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// ASHTTPQuery 表示解析后的ActiveSync HTTP查询结构
type ASHTTPQuery struct {
	ProtocolVersion uint8   // 协议版本 (1字节)
	CommandCode     uint8   // 命令代码 (1字节)
	Locale          uint16  // 区域设置 (2字节)
	DeviceID        string  // 设备ID
	PolicyKey       *uint32 // 策略密钥 (可选，4字节)
	DeviceType      string  // 设备类型
	CommandParams   []byte  // 命令参数
}

// Command codes
const (
	CmdSync              uint8 = 0  // 同步文件夹变更
	CmdSendMail          uint8 = 1  // 发送邮件
	CmdSmartForward      uint8 = 2  // 转发邮件
	CmdSmartReply        uint8 = 3  // 回复邮件
	CmdGetAttachment     uint8 = 4  // 获取附件
	CmdFolderSync        uint8 = 9  // 同步文件夹层级
	CmdFolderCreate      uint8 = 10 // 创建文件夹
	CmdFolderDelete      uint8 = 11 // 删除文件夹
	CmdFolderUpdate      uint8 = 12 // 移动或重命名文件夹
	CmdMoveItems         uint8 = 13 // 移动项目
	CmdGetItemEstimate   uint8 = 14 // 获取同步项目数量估计
	CmdMeetingResponse   uint8 = 15 // 处理会议请求
	CmdSearch            uint8 = 16 // 搜索全局通讯录
	CmdSettings          uint8 = 17 // 获取和设置全局属性
	CmdPing              uint8 = 18 // 监控文件夹变更
	CmdItemOperations    uint8 = 19 // 项目操作
	CmdProvision         uint8 = 20 // 获取安全策略设置
	CmdResolveRecipients uint8 = 21 // 解析收件人
)

// commandNames 提供命令代码到名称的映射
// 参考文档：https://learn.microsoft.com/en-us/openspecs/exchange_server_protocols/ms-ashttp/0ab55ebc-6ea9-4ae4-af37-5736d5195d46
var commandNames = map[uint8]string{
	CmdSync:              "Sync",
	CmdSendMail:          "SendMail",
	CmdSmartForward:      "SmartForward",
	CmdSmartReply:        "SmartReply",
	CmdGetAttachment:     "GetAttachment",
	CmdFolderSync:        "FolderSync",
	CmdFolderCreate:      "FolderCreate",
	CmdFolderDelete:      "FolderDelete",
	CmdFolderUpdate:      "FolderUpdate",
	CmdMoveItems:         "MoveItems",
	CmdGetItemEstimate:   "GetItemEstimate",
	CmdMeetingResponse:   "MeetingResponse",
	CmdSearch:            "Search",
	CmdSettings:          "Settings",
	CmdPing:              "Ping",
	CmdItemOperations:    "ItemOperations",
	CmdProvision:         "Provision",
	CmdResolveRecipients: "ResolveRecipients",
}

// ParseBase64Query 解析base64编码的查询字符串
// 参考文档：https://learn.microsoft.com/en-us/openspecs/exchange_server_protocols/ms-ashttp/9f75a516-edff-48d2-94b9-5f2a3cc570d3
func ParseBase64Query(base64Query string) (*ASHTTPQuery, error) {
	// 解码base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Query)
	if err != nil {
		return nil, err
	}

	if len(data) < 5 { // 最小长度检查：版本(1) + 命令码(1) + 区域(2) + 设备ID长度(1)
		return nil, errors.New("query data too short")
	}

	query := &ASHTTPQuery{}

	// 读取固定长度字段
	offset := 0
	query.ProtocolVersion = data[offset]
	offset++

	query.CommandCode = data[offset]
	offset++

	query.Locale = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	// 读取设备ID
	if offset >= len(data) {
		return nil, errors.New("unexpected end of data at device ID length")
	}
	deviceIDLen := int(data[offset])
	offset++

	if deviceIDLen <= 0 {
		return nil, errors.New("invalid device ID length")
	}
	if offset+deviceIDLen > len(data) {
		return nil, errors.New("unexpected end of data at device ID")
	}
	query.DeviceID = string(data[offset : offset+deviceIDLen])
	offset += deviceIDLen

	// 读取策略密钥（如果存在）
	if offset < len(data) {
		policyKeyLen := int(data[offset])
		offset++

		if policyKeyLen == 4 && offset+4 <= len(data) {
			policyKey := binary.LittleEndian.Uint32(data[offset : offset+4])
			query.PolicyKey = &policyKey
			offset += 4
		} else if policyKeyLen != 0 {
			return nil, errors.New("invalid policy key length")
		}
	}

	// 读取设备类型（如果存在）
	if offset < len(data) {
		deviceTypeLen := int(data[offset])
		offset++

		if offset+deviceTypeLen <= len(data) {
			query.DeviceType = string(data[offset : offset+deviceTypeLen])
			offset += deviceTypeLen
		} else {
			return nil, errors.New("unexpected end of data at device type")
		}
	}

	// 剩余的数据作为命令参数
	if offset < len(data) {
		query.CommandParams = data[offset:]
	}

	return query, nil
}

// GetCommandName 返回命令代码对应的名称
func (q ASHTTPQuery) GetCommandName() string {
	if name, ok := commandNames[q.CommandCode]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(%d)", q.CommandCode)
}

// MarshalJSON 自定义JSON序列化方法
func (q ASHTTPQuery) MarshalJSON() ([]byte, error) {
	type Alias ASHTTPQuery
	return json.Marshal(&struct {
		DeviceID    string `json:"DeviceID"`
		CommandName string `json:"CommandName"`
		*Alias
	}{
		DeviceID:    hex.EncodeToString([]byte(q.DeviceID)),
		CommandName: q.GetCommandName(),
		Alias:       (*Alias)(&q),
	})
}
