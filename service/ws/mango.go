package ws

import (
	"gin-IM/db/mango"
	"gin-IM/model/ws"
	"time"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

// InsertMsg 向MangoDB中插入信息
func InsertMsg(database string, id string, content string, read bool, expire int64) error {
	// 获取collection
	collection := mango.MgClient.Database(database).Collection(id) // 设计上，集合名为1-->2之类的（即传过来的id）
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err := collection.InsertOne(mango.MgCtx, comment)
	if err != nil {
		return err
	}
	return nil
}
