package nest

// 参数
type MsgParams struct {
	Msg  interface{} `title:"消息内容"`
	Sign string      `title:"认证签名"`
	Code string      `title:"消息码"`
}

// 返回数据
type MsgResult struct {
	Id       int64       `xorm:"not null pk autoincr BIGINT(20)" title:"消息ID"`
	TmpId    string      `title:"临时消息uuid"`
	Platform int         `title:"消息平台 0表示所有平台"`
	UserId   int64       `title:"接收者ID" check:"notnull"`
	Code     string      `title:"消息码" check:"notnull"`
	Msg      interface{} `title:"消息内容" check:"notnull"`
}
