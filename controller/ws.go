package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/luxingwen/secret-game/dao"
	"net/http"
	"time"
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
	t := time.NewTicker(time.Minute)
	for {
		select {
		case c := <-WsManager.Register:
			WsManager.MWsClient[c.Id] = c

		case c := <-WsManager.Unregister:
			delete(WsManager.MWsClient, c.Id)
		case <-t.C:
			msg := Message{Cmd: "servertime", Data: time.Now().Unix()}
			b, err := json.Marshal(msg)
			if err != nil {
				fmt.Println("json marshl err:", err)
				continue
			}
			for _, c := range manager.MWsClient {
				c.Send <- b
			}
		}
	}
}

func (c *WsClient) Read() {
	defer func() {
		WsManager.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("err:", err)
			WsManager.Unregister <- c
			c.Conn.Close()
			break
		}
		fmt.Println("rev:", string(msg))
		msgData := new(Message)
		err = json.Unmarshal(msg, &msgData)
		if err != nil {
			fmt.Println("json unmarshl er:", err)
			continue
		}
		c.handlerMsg(msgData)

	}

}

func (c *WsClient) handlerMsg(msg *Message) {
	if msg.Cmd == "register" {
		token := msg.Data.(string)
		claims, err := ParseWxToken(token)
		if err != nil {
			fmt.Println("parse wx token err:", err)
			return
		}

		c.Id = claims.Id
		WsManager.Register <- c
		fmt.Println("c-->", claims)
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
	fmt.Println("--->ws")
	uid := c.GetInt("wxUserId")
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("2222->")
		http.NotFound(c.Writer, c.Request)
		return
	}

	fmt.Println("ws--id>", uid)
	client := &WsClient{
		Conn: conn,
		Send: make(chan []byte),
	}

	if uid != 0 {
		client.Id = uid
		WsManager.Register <- client

	} else {
		fmt.Println("无效的客户端-->", uid)
	}

	go client.Read()
	go client.Write()

}

func NotifyTeams(uid int, cmd string, data interface{}) {
	msg := Message{Cmd: cmd, Data: data}
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json marshl err:", err)
		return
	}

	teamUsers, err := dao.GetDao().GetTeamUserMapsByUid(uid)
	if err != nil {
		fmt.Println("GetTeamUserMapsByUid err:", err)
		return
	}

	for _, item := range teamUsers {
		userId := int(item.UserId)
		if userId == uid {
			continue
		}
		if c, ok := WsManager.MWsClient[userId]; ok {
			c.Send <- b
		}
	}
}
