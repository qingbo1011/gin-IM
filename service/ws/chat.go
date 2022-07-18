package ws

import (
	"encoding/json"
	"fmt"
	"gin-IM/db/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

const month = 60 * 60 * 24 * 30 // 一个月30天

// SendMsg 发送数据的类型
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// ReplyMsg 回复的消息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// Client 用户类
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

func (c *Client) Read() {
	defer func() { // 别忘了Close
		Manager.Unregister <- c
		err := c.Socket.Close()
		if err != nil {
			logging.Info(err)
		}
	}()
	for {
		c.Socket.PongHandler()
		//sendMsg := new(SendMsg)
		var sendMsg *SendMsg
		// _,msg,_:=c.Socket.ReadMessage()	// 传来的不是json
		err := c.Socket.ReadJSON(&sendMsg)
		if err != nil {
			logging.Info(err)
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		switch sendMsg.Type {
		// 发送消息
		case 1:
			// 这个功能我感觉可以不用加。但是以学习的态度，还是写一下
			// 限制单向骚扰和轰炸
			r1, err := redis.Rdb.Get(redis.RCtx, c.ID).Result()
			r2, err := redis.Rdb.Get(redis.RCtx, c.SendID).Result()
			if err != nil {
				logging.Info(err)
			}
			// 发了3条信息，对方1个没回，此时限制发送
			if r1 >= "3" && r2 == "" { // go语言支持 "3">"2"这种字符串比较
				replayMsg := ReplyMsg{
					Code:    http.StatusForbidden,
					Content: "到达限制，请等待对方回复",
				}
				msg, err := json.Marshal(replayMsg)
				if err != nil {
					logging.Fatalln(err)
				}
				err = c.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					logging.Info(err)
				}
				redis.Rdb.Expire(redis.RCtx, c.ID, time.Hour*24*30) // 设置30天的过期时间
				continue
			} else {
				redis.Rdb.Incr(redis.RCtx, c.ID)
				redis.Rdb.Expire(redis.RCtx, c.ID, time.Hour*24*30*3) // 设置90天的过期时间
			}
			fmt.Println(c.ID, "发送信息", sendMsg.Content)
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		// 拉取历史消息
		case 2:
			t, err := strconv.Atoi(sendMsg.Content) // 传送来时间
			if err != nil {
				t = 999999999
			}
			fmt.Println(t)
		//
		case 3:

		}

	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if ok {
				err := c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logging.Info(err)
				}
				return
			}
			fmt.Println(c.ID, "接受消息", string(message))
			replyMsg := ReplyMsg{
				Code:    http.StatusOK,
				Content: string(message),
			}
			msg, err := json.Marshal(replyMsg)
			if err != nil {
				logging.Fatal(err)
			}
			err = c.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logging.Info(err)
			}
		}
	}
}

// Broadcast 广播类，包括广播内容和广播源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ClientManager 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Message 信息转JSON（包括：发送者、接受者、内容）
type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

var Manager = ClientManager{ // 全局的用户管理Manger
	Clients:    make(map[string]*Client), // 参与连接的用户（出于性能的考虑，需要设置最大连接数）
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func WsHandler(c *gin.Context) {
	uid := c.Query("uid")     // 自己的id
	toUid := c.Query("toUid") // 对方的id
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建一个用户实例
	client := &Client{
		ID:     createId(uid, toUid),
		SendID: createId(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理Manager上
	Manager.Register <- client
	go client.Read()
	go client.Write()

}

func createId(uid, toUid string) string {
	return uid + "-->" + toUid
}
