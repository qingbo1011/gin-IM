package ws

import (
	"fmt"
	"gin-IM/db/mango"
	"gin-IM/model/ws"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     bool   `json:"read"`
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

func FindMany(db string, sendId string, id string, time int64, pageSize int64) ([]ws.Result, error) {
	var resultMe []ws.Trainer
	var resultYou []ws.Trainer
	// 获取collection
	sendIdCollection := mango.MgClient.Database(db).Collection(sendId)
	idCollection := mango.MgClient.Database(db).Collection(id)
	// SetSort 设置排序字段（1表示升序；-1表示降序）
	op1 := options.Find().SetSort(bson.D{{"starttime", -1}})
	op2 := options.Find().SetLimit(int64(pageSize))
	sendIdcur, err := sendIdCollection.Find(mango.MgCtx, op1, op2)
	idcur, err := idCollection.Find(mango.MgCtx, op1, op2)
	err = sendIdcur.All(mango.MgCtx, &resultYou) // sendId 对面发过来的
	err = idcur.All(mango.MgCtx, &resultMe)      // id 发给对面的
	if err != nil {
		return nil, err
	}
	results, err := appendAndSort(resultMe, resultYou)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func FirstFindMsg(db string, sendId string, id string) ([]ws.Result, error) {
	var results []ws.Result
	// 首次查询(把对方发来的所有未读都取出来)
	var resultMe []ws.Trainer
	var resultYou []ws.Trainer
	// 获取collection
	sendIdCollection := mango.MgClient.Database(db).Collection(sendId)
	idCollection := mango.MgClient.Database(db).Collection(id)
	filter := bson.M{
		"read": bson.M{
			"&all": []bool{false},
		},
	}
	op1 := options.Find().SetSort(bson.D{{"starttime", 1}})
	op2 := options.Find().SetLimit(1)
	sendCur, err := sendIdCollection.Find(mango.MgCtx, filter, op1, op2)
	if sendCur == nil {
		return nil, err
	}
	var unReads []ws.Trainer
	err = sendCur.All(mango.MgCtx, &unReads)
	if err != nil {
		return nil, err
	}
	if len(unReads) > 0 {
		timeFilter := bson.M{
			"starttime": bson.M{
				"$gte": unReads[0].StartTime,
			},
		}
		sendIdTimeCur, err := sendIdCollection.Find(mango.MgCtx, timeFilter)
		idTimeCur, err := idCollection.Find(mango.MgCtx, timeFilter)
		err = sendIdTimeCur.All(mango.MgCtx, &resultYou)
		err = idTimeCur.All(mango.MgCtx, &resultMe)
		if err != nil {
			return nil, err
		}
		results, err := appendAndSort(resultMe, resultYou)
		if err != nil {
			return nil, err
		}
		return results, nil
	} else {
		results, err = FindMany(db, sendId, id, 99999999, 10)
	}
	overTime := bson.D{
		{
			"$and", bson.A{
				bson.D{{"endtime", bson.M{"&lt": time.Now().Unix()}}},
				bson.D{{"read", bson.M{"$eq": true}}},
			}},
	}
	// 删除过期信息
	_, err = sendIdCollection.DeleteMany(mango.MgCtx, overTime)
	_, err = idCollection.DeleteMany(mango.MgCtx, overTime)
	// 将所有的维度设置已读
	_, err = sendIdCollection.UpdateMany(mango.MgCtx, filter, bson.M{
		"$set": bson.M{"read": true},
	})
	_, err = sendIdCollection.UpdateMany(mango.MgCtx, filter, bson.M{
		"&set": bson.M{"endtime": time.Now().Unix() + int64(month*3)},
	})
	return results, nil
}

func appendAndSort(resultMe []ws.Trainer, resultYou []ws.Trainer) ([]ws.Result, error) {
	var results []ws.Result
	for _, v := range resultMe {
		sendSort := SendSortMsg{
			Content:  v.Content,
			Read:     v.Read,
			CreateAt: v.StartTime,
		}
		result := ws.Result{
			StartTime: v.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, v := range resultYou {
		sendSort := SendSortMsg{
			Content:  v.Content,
			Read:     v.Read,
			CreateAt: v.StartTime,
		}
		result := ws.Result{
			StartTime: v.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	// 最后进行排序，按StartTime升序排序（自然排序）
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime < results[j].StartTime
	})
	return results, nil
}
