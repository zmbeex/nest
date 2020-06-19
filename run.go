package nest

import (
	"github.com/gorilla/websocket"
	"github.com/zmbeex/gkit"
	"net/http"
)

var ug = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LoginNestParams struct {
	Token    string // token
	World    string // 所属世界
	Time     int64  // 时间戳
	Platform int    // 登录平台
}

func Run() {
	http.HandleFunc("/nest", func(w http.ResponseWriter, r *http.Request) {
		conn, err := ug.Upgrade(w, r, nil)
		gkit.CheckPanic(err, "websocket连接失败")

		//sign := r.FormValue("sign")
		//loginParams := new(LoginNestParams)
		//token, err := ttoken.GetToken(sign)
		//if err != nil {
		//	_ = conn.WriteJSON(map[string]interface{}{
		//		"note": "请先登录",
		//		"code": -1,
		//		"info": err.Error(),
		//	})
		//	_ = conn.Close()
		//	return
		//}

		// 创建客户端
		c := &Client{
			Conn: conn,
			Send: make(chan *MsgResult, Cache.Set.ClientChanLen),
			//UserId:   token.Id,
			//Platform: loginParams.Platform,
		}
		c.InitWork()

		go c.ReadMessage()
		go c.SendMessage()
	})
}
