package nest

import (
	"errors"
	"github.com/zmbeex/gkit"
	"go.uber.org/zap"
	"sync"
)

type User struct {
	Id     int64 `title:"用户ID"`
	Client sync.Map
}

// 获取user对象
func GetUser(id int64) *User {
	u, ok := Cache.UserPool.Load(id)
	if ok {
		user, ok := u.(*User)
		if ok {
			return user
		}
	}
	user := new(User)
	user.Id = id
	Cache.UserPool.Store(id, user)
	return user
}

// 增加客户端
func (m *User) AddClient(c *Client) {
	m.Client.Store(c.Platform, c)
	gkit.Debug("[NEST]增加客户端", zap.Int64("userId", m.Id), zap.Int("platform", c.Platform))
}

// 删除客户端
func (m *User) DeleteClient(platform int) {
	m.Client.Delete(platform)
	gkit.Debug("[NEST]删除客户端", zap.Int64("userId", m.Id), zap.Int("platform", platform))
}

// 销毁用户，当客户端为空时
func (m *User) destroyWhenClientNull() {
	isExistClient := true
	m.Client.Range(func(key, value interface{}) bool {
		isExistClient = false
		return false
	})
	if !isExistClient {
		Cache.UserPool.Delete(m.Id)
		gkit.Debug("[NEST]从用户池删除用户", zap.Int64("userId", m.Id))
	}
}

// 获取所有客户端
func (m *User) GetAllClient() map[int]*Client {
	data := make(map[int]*Client)
	m.Client.Range(func(key, value interface{}) bool {
		platform := gkit.ToInt(key)
		if platform > 0 {
			c, ok := value.(*Client)
			if ok {
				data[platform] = c
			}
		}
		return true
	})
	return data
}

// 获取某个平台的客户端
func (m *User) GetClient(platform int) *Client {
	val, ok := m.Client.Load(platform)
	if !ok {
		return nil
	}
	c, ok := val.(*Client)
	if ok {
		return c
	}
	return nil
}

// 消息发送前处理
func (m *User) sendBeforeHandle(msg *MsgResult) error {
	// 消息码不能为空
	if msg.Code == "" {
		return errors.New("消息码不能为空")
	}
	// 消息接收者不能为空
	if msg.UserId == 0 {
		msg.UserId = m.Id
	}
	// 判断消息是否需要持久化
	if msg.TmpId == "" {
		// 持久消息
		Handler.PersistenceSaveUserMsg(msg)
	}
	return nil
}

// 发送给用户的所有客户端
func (m *User) Send(msg *MsgResult) {
	err := m.sendBeforeHandle(msg)
	gkit.CheckPanic(err, "消息持久化失败")

	allClient := m.GetAllClient()
	for _, c := range allClient {
		c.Send <- msg
	}
	gkit.Debug("[NEST]发送给用户的所有客户端", zap.Int64("userId", m.Id))
}

// 发送给某个平台
func (m *User) SendPlatform(platform int, msg *MsgResult) {
	err := m.sendBeforeHandle(msg)
	gkit.CheckPanic(err, "消息持久化失败")

	c := m.GetClient(platform)
	c.Send <- msg
	gkit.Debug("[NEST]发送给某个平台", zap.Int64("userId", m.Id), zap.Int("platform", c.Platform))
}
