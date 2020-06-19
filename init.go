package nest

import (
	"errors"
	"github.com/zmbeex/gkit"
	"sync"
)

type Router struct {
	Code   string
	Params interface{}
	Desc   string
	Fn     func(n *Nest)
}

type Set struct {
	// 客户端发送消息信道长度
	ClientChanLen   int    `title:"客户端发送消息信道长度"`
	LoginSignAesKey string `title:"登录加密私钥"`
}

// 数据缓存
var Cache struct {
	Set *Set

	ClientPool sync.Map `title:"客户端连接池 map[*Client]boo"`
	UserPool   sync.Map `title:"用户的客户端 map[int64]*User"`
	UserRoom   sync.Map `title:"用户级房间  map[int64]*User"`
	ClientRoom sync.Map `title:"客户端级别的房间  map[int64]*Client"`

	// 路由
	Router map[string]*Router
	// 中间件列表
	MiddlewareList []func(n *Nest)
}

// 配置方法
var Handler struct {
	// 持久化消息
	PersistenceSaveUserMsg func(msg *MsgResult)
}

// 初始化
func InitSelf(set *Set) {
	Cache.Set = set
	if Cache.Set.ClientChanLen == 0 {
		Cache.Set.ClientChanLen = 100
		gkit.Info("[NEST]请设置nest客户端发送消息信道长度，默认100")
	}
	// 初始化路由存储
	Cache.Router = make(map[string]*Router)

	// 持久化消息
	Handler.PersistenceSaveUserMsg = func(msg *MsgResult) {
		gkit.Warn("请配置持久化: nest.Handler.PersistenceSaveUserMsg")
	}
}

// 注册路由
func Register(code string, desc string, params interface{}, fn func(n *Nest)) {
	route := new(Router)
	route.Code = code
	route.Desc = desc
	route.Params = params
	route.Fn = fn
	Cache.Router[code] = route
}

// 获取可执行方法
func GetExecFn(key string) (error, func(n *Nest)) {
	route := Cache.Router[key]
	if route == nil {
		return errors.New("路由不存在: " + key), nil
	}
	return nil, route.Fn
}
