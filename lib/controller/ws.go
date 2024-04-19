package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wike2019/wike_go/lib/utils"
	"net/http"
	"sync"
	"time"
)

const Default = "default"

// websocket 广播函数
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket 注册函数
func WsHandler(c *gin.Context) {
	roomId := c.DefaultQuery("room", Default)
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	record := WC.Data.Get(roomId)
	record = append(record, conn)
	WC.Data.Set(roomId, record)
}

var WC *WsClient

func init() {
	WC = &WsClient{
		Data: utils.NewMap[[]*websocket.Conn](),
	}
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			for _, k := range WC.Data.Keys() {
				ping(k)
			}
		}
	}()
}

type WsClient struct {
	Data *utils.MapSync[[]*websocket.Conn]
	lock sync.Mutex
}

func Broadcast(room string, data string) {
	// 创建一个新的切片存储仍然有效的连接
	activeConns := make([]*websocket.Conn, 0)
	connList := WC.Data.Get(room)
	for _, conn := range connList {
		conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
		// 尝试向连接发送消息
		if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
			conn.Close()
		} else {
			activeConns = append(activeConns, conn)
		}
	}

	WC.Data.Set(room, activeConns)
}
func ping(room string) {
	// 创建一个新的切片存储仍然有效的连接
	activeConns := make([]*websocket.Conn, 0)
	connList := WC.Data.Get(room)
	for _, conn := range connList {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		// 尝试向连接发送消息
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			conn.Close()
		} else {
			activeConns = append(activeConns, conn)
		}
	}
	WC.Data.Set(room, activeConns)
}
