package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type WsClient struct {
	Id   int
	Conn *websocket.Conn
	Send chan []byte
}

type WsClientManager struct {
	MWsClient  map[int]*WsClient
	Register   chan *WsClient
	Unregister chan *WsClient
}

type Message struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

var WsManager = &WsClientManager{
	MWsClient:  make(map[int]*WsClient, 0),
	Register:   make(chan *WsClient),
	Unregister: make(chan *WsClient),
}

func (manager *WsClientManager) Start() {
	for {
		select {
		case c := <-WsManager.Register:
			WsManager.MWsClient[c.Id] = c

		case c := <-WsManager.Unregister:
			delete(WsManager.MWsClient, c.Id)
		}
	}
}

func (c *WsClient) Read() {
	defer func() {
		WsManager.Unregister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PongHandler()
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			WsManager.Unregister <- c
			c.Conn.Close()
			break
		}
		_ = msg
	}

}

func (c *WsClient) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func WsHandler(c *gin.Context) {
	uid := c.GetInt("wxUserId")
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	client := &WsClient{
		Id:   uid,
		Conn: conn,
		Send: make(chan []byte),
	}
	WsManager.Register <- client
	go client.Read()
	go client.Write()

}
