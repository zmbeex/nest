package nest

import "github.com/zmbeex/gkit"

type Middleware struct {
	Index    int
	Handlers []func(n *Nest)
}

// 加入中间件
func Use(f func(nest *Nest)) {
	Cache.MiddlewareList = append(Cache.MiddlewareList, f)
}

// 中间件 执行
func (c *Nest) Next() {
	c.Middleware.Index++
	for s := len(c.Middleware.Handlers); c.Middleware.Index < s; c.Middleware.Index++ {
		c.Middleware.Handlers[c.Middleware.Index](c)
	}
}

// 关闭中间件
func (c *Nest) Close() {
	c.Middleware.Index = 999999999
}

// 执行中间件
func (c *Nest) runMiddleware(f func(n *Nest)) {
	c.Middleware.Handlers = append(Cache.MiddlewareList, f)
	c.Middleware.Index = -1
	c.Next()
}

// 标准异常处理
func CatchHandleMiddleware() func(n *Nest) {
	return func(n *Nest) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			// 普通异常
			err, ok := r.(error)
			if ok {
				n.ErrorMsg(err.Error(), "CatchHandleMiddleware.error")
				return
			}
			// 自定义异常
			err2, ok := r.(*gkit.Zerror)
			if ok {
				n.ErrorMsg(err2.Note, "CatchHandleMiddleware.Zerror")
				return
			}
			// 未知异常
			n.ErrorMsg("未知异常", "CatchHandleMiddleware.unknown")
		}()
		// 执行其他操作
		n.Next()
	}
}
