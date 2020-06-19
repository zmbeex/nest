package nest

type Nest struct {
	UserId     int64      `title:"用户ID"`
	User       *User      `title:"用户连接对象"`
	Params     *MsgParams `title:"消息参数"`
	Result     *MsgResult `title:"返回结果"`
	Middleware *Middleware
}

func (c *Nest) ErrorMsg(msg string, code string) {
	// todo 处理返回异常数据

}
