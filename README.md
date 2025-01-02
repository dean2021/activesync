# ActiveSync HTTP Query Parser

这是一个用于解析Microsoft Exchange ActiveSync HTTP查询的Go语言库。它可以解析Base64编码的查询字符串，提取出其中包含的各种协议字段。

## 功能特性

- 解析Base64编码的ActiveSync HTTP查询字符串
- 支持所有标准ActiveSync命令代码
- 提供友好的JSON输出格式
- 设备ID自动转换为十六进制显示
- 命令代码自动映射为可读名称

## 安装

```bash
go get github.com/dean2021/activesync
```

## 使用示例

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/dean2021/activesync/ashttp"
)

func main() {
    // Base64编码的查询字符串
    query, err := ashttp.ParseBase64Query("oQkECBCeDEK6NjuTWKLjgUH2WCxdBIIanKgLV2luZG93c01haWw=")
    if err != nil {
        log.Fatal(err)
    }

    // 格式化输出
    b, err := json.MarshalIndent(query, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(b))
}
```

输出示例：
```json
{
  "DeviceID": "9e0c42ba363b9358a2e38141f6582c5d",
  "CommandName": "FolderSync",
  "ProtocolVersion": 161,
  "CommandCode": 9,
  "Locale": 2052,
  "PolicyKey": 2768931330,
  "DeviceType": "WindowsMail",
  "CommandParams": null
}
```

## 支持的命令

| 命令代码 | 命令名称 | 描述 |
|---------|---------|------|
| 0 | Sync | 同步文件夹变更 |
| 1 | SendMail | 发送邮件 |
| 2 | SmartForward | 转发邮件 |
| 3 | SmartReply | 回复邮件 |
| 4 | GetAttachment | 获取附件 |
| 9 | FolderSync | 同步文件夹层级 |
| 10 | FolderCreate | 创建文件夹 |
| 11 | FolderDelete | 删除文件夹 |
| 12 | FolderUpdate | 移动或重命名文件夹 |
| 13 | MoveItems | 移动项目 |
| 14 | GetItemEstimate | 获取同步项目数量估计 |
| 15 | MeetingResponse | 处理会议请求 |
| 16 | Search | 搜索全局通讯录 |
| 17 | Settings | 获取和设置全局属性 |
| 18 | Ping | 监控文件夹变更 |
| 19 | ItemOperations | 项目操作 |
| 20 | Provision | 获取安全策略设置 |
| 21 | ResolveRecipients | 解析收件人 |

## 数据结构

查询结构体定义：
```go
type ASHTTPQuery struct {
    ProtocolVersion uint8   // 协议版本 (1字节)
    CommandCode     uint8   // 命令代码 (1字节)
    Locale          uint16  // 区域设置 (2字节)
    DeviceID        string  // 设备ID
    PolicyKey       *uint32 // 策略密钥 (可选，4字节)
    DeviceType      string  // 设备类型
    CommandParams   []byte  // 命令参数
}
```

## 参考文档

- [MS-ASHTTP: Exchange ActiveSync: HTTP Protocol](https://learn.microsoft.com/en-us/openspecs/exchange_server_protocols/ms-ashttp/9f75a516-edff-48d2-94b9-5f2a3cc570d3)
- [Command Codes](https://learn.microsoft.com/en-us/openspecs/exchange_server_protocols/ms-ashttp/0ab55ebc-6ea9-4ae4-af37-5736d5195d46)

## 许可证

MIT License
