package ashttp

import (
	"encoding/base64"
	"testing"
)

func TestParseBase64Query(t *testing.T) {
	// 创建一个测试用的查询数据
	testData := []byte{
		141,                // 协议版本 (141 = 0x8D)
		1,                  // 命令代码
		0x09, 0x04,        // 区域设置 (1033 = 英语)
		5,                  // 设备ID长度
		't', 'e', 's', 't', '1', // 设备ID
		4,                  // 策略密钥长度
		1, 0, 0, 0,        // 策略密钥值 (1)
		5,                  // 设备类型长度
		'p', 'h', 'o', 'n', 'e', // 设备类型
		'p', 'a', 'r', 'a', 'm', 's', // 命令参数
	}

	// 将测试数据转换为base64
	base64Query := base64.StdEncoding.EncodeToString(testData)

	// 解析数据
	query, err := ParseBase64Query(base64Query)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	// 验证解析结果
	if query.ProtocolVersion != 141 {
		t.Errorf("协议版本不匹配, 期望 141, 得到 %d", query.ProtocolVersion)
	}

	if query.CommandCode != 1 {
		t.Errorf("命令代码不匹配, 期望 1, 得到 %d", query.CommandCode)
	}

	if query.Locale != 1033 {
		t.Errorf("区域设置不匹配, 期望 1033, 得到 %d", query.Locale)
	}

	if query.DeviceID != "test1" {
		t.Errorf("设备ID不匹配, 期望 'test1', 得到 '%s'", query.DeviceID)
	}

	if query.PolicyKey == nil || *query.PolicyKey != 1 {
		t.Errorf("策略密钥不匹配, 期望 1, 得到 %v", query.PolicyKey)
	}

	if query.DeviceType != "phone" {
		t.Errorf("设备类型不匹配, 期望 'phone', 得到 '%s'", query.DeviceType)
	}

	expectedParams := []byte("params")
	if string(query.CommandParams) != string(expectedParams) {
		t.Errorf("命令参数不匹配, 期望 '%s', 得到 '%s'", expectedParams, query.CommandParams)
	}
}
