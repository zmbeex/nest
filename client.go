package nest

import (
	"github.com/gorilla/websocket"
	"github.com/zmbeex/gkit"
	"go.uber.org/zap"
)

type Client struct {
	UserId   int64           `title:"用户ID"`
	Platform int             `title:"平台"`
	Conn     *websocket.Conn `title:"连接"`
	Send     chan *MsgResult `title:"发送消息"`
}

// 初始化客户端
func (m *Client) InitWork() {
	// 用户加入用户池
	val, ok := Cache.UserPool.Load(m.UserId)
	user, ok := val.(*User)
	if !ok {
		user = new(User)
		user.Id = m.UserId
		Cache.UserPool.Store(user.Id, user)
		gkit.Debug("[NEST]用户加入用户池", zap.Int64("userId", user.Id))
	}

	// 客户端加入连接池
	Cache.ClientPool.Store(m, true)
	gkit.Debug("[NEST]客户端加入连接池", zap.Int64("userId", m.UserId), zap.Int("platform", m.Platform))

	// 客户端加入到用户中
	user.AddClient(m)
	gkit.Debug("[NEST]客户端加入到用户客户池", zap.Int64("userId", m.UserId), zap.Int("platform", m.Platform))

}

// 客户端销毁
func (m *Client) destroy() {
	// 客户端加入连接池
	Cache.ClientPool.Delete(m)

	// 客户加入用户组
	val, ok := Cache.UserPool.Load(m.UserId)
	user, ok := val.(*User)
	if !ok {
		gkit.Error("[NEST]客户端没有有效从属: ", zap.Int64("userId", m.UserId), zap.Any("user", user))
		return
	}

	// 删除用户的某个客户端
	user.DeleteClient(m.Platform)
	gkit.Debug("[NEST]删除用户的某个客户端", zap.Int64("userId", m.UserId), zap.Int("platform", m.Platform))

	// 判断用户是否没有客户端
	user.destroyWhenClientNull()

}

// 读取消息
func (c *Client) ReadMessage() {
	for {
		// 接收消息
		message := &MsgParams{}
		if err := c.Conn.ReadJSON(message); err != nil {
			gkit.Info("[NEST]消息发送异常（自动断开连接）: " + err.Error())
			// 销毁客户端
			c.destroy()
			return
		}
		gkit.Info("[NEST]读取消息", zap.Any("ReadMessage", message))
		nest := new(Nest)
		// 返回消息
		nest.Result = new(MsgResult)
		nest.Result.Code = message.Code
		// 来源入参
		nest.Params = new(MsgParams)
		*nest.Params = *message
		// 用户ID
		nest.UserId = c.UserId
		// 用户对象
		nest.User = GetUser(c.UserId)
		// 中间件
		nest.Middleware = new(Middleware)

		go func() {
			// 进程保护
			defer gkit.DeferRecover("[NEST]未知异常")
			// todo 处理数据
			err, fn := GetExecFn(nest.Params.Code)
			gkit.CheckErrLog(err, "无效Code")
			if fn != nil {
				// 执行处理方法
				nest.runMiddleware(fn)
			} else {
				gkit.Info("[NEST]没有注册方法:" + nest.Params.Code)
			}
		}()
	}
}

/// 发送消息
func (c *Client) SendMessage() {
	for {
		// 接收信道的消息
		result := <-c.Send
		gkit.Info("发送给客户端: " + gkit.SetJson(result))
		// 发送给客户端
		if err := c.Conn.WriteJSON(result); err != nil {
			// 销毁客户端
			c.destroy()
			return
		}
	}
}
