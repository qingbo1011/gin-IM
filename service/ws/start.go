package ws

import (
	"encoding/json"
	"fmt"
	config "gin-IM/conf"
	"net/http"

	"github.com/gorilla/websocket"
	logging "github.com/sirupsen/logrus"
)

func (m *ClientManager) Start() {
	for {
		fmt.Println("<---监听管道通信--->")
		select {
		// 建立连接
		case conn := <-Manager.Register:
			fmt.Println("建立新连接: ", conn.ID)
			Manager.Clients[conn.ID] = conn
			replayMsg := &ReplyMsg{
				Code:    http.StatusOK,
				Content: "已连接至服务器",
			}
			msg, err := json.Marshal(replayMsg)
			if err != nil {
				logging.Info(err)
			}
			err = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logging.Info(err)
			}
		// 断开连接
		case conn := <-Manager.Unregister:
			fmt.Println("连接失败: ", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replayMsg := &ReplyMsg{
					Code:    http.StatusInternalServerError,
					Content: "连接已断开",
				}
				msg, err := json.Marshal(replayMsg)
				if err != nil {
					logging.Info(err)
				}
				err = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					logging.Info(err)
				}
				close(conn.Send)
				delete(Manager.Clients, conn.ID) // 使用delete()函数从map中删除一组键值对
			}
		// 广播信息
		case broadcast := <-Manager.Broadcast:
			message := broadcast.Message
			sendID := broadcast.Client.SendID
			flag := false // 默认对方不在线
			for id, conn := range Manager.Clients {
				if id != sendID {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID
			if flag {
				fmt.Println("对方在线应答")
				replayMsg := ReplyMsg{
					Code:    http.StatusOK,
					Content: "对方在线应答",
				}
				msg, err := json.Marshal(replayMsg)
				if err != nil {
					logging.Info(err)
				}
				err = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					logging.Info(err)
				}
				// 将消息插入到MangoDB中
				// true表示已读（这里处理逻辑比较粗糙，在线就认为已读）
				err = InsertMsg(config.MangoDBName, id, string(message), true, int64(month*3))
				if err != nil {
					logging.Info(err)
				}
			} else {
				fmt.Println("对方不在线")
				replayMsg := ReplyMsg{
					Code:    http.StatusNotFound,
					Content: "对方不在线",
				}
				msg, err := json.Marshal(replayMsg)
				if err != nil {
					logging.Info(err)
				}
				err = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					logging.Info(err)
				}
				err = InsertMsg(config.MangoDBName, id, string(message), false, int64(month*3))
				if err != nil {
					logging.Info(err)
				}
			}
		}
	}
}
