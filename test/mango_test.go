package test

import (
	"fmt"
	config "gin-IM/conf"
	"gin-IM/db/mango"
	"gin-IM/model/ws"
	"sort"
	"testing"

	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	config.Init("../conf/config.ini")
	mango.Init()
}

func TestMango(t *testing.T) {
	collection := mango.MgClient.Database("test").Collection("person")
	// Find
	// SetSort 设置排序字段（1表示升序；-1表示降序）
	findOptions := options.Find().SetSort(bson.D{{"level", 1}})
	findCursor, err := collection.Find(mango.MgCtx, bson.M{"gender": "男"}, findOptions)
	if err != nil {
		logging.Fatal(err)
	}
	var results []bson.M
	err = findCursor.All(mango.MgCtx, &results)
	if err != nil {
		logging.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func TestBoolFind(t *testing.T) {
	collection := mango.MgClient.Database("gin-IM").Collection("1-->2")
	filter := bson.M{
		"read": bson.M{
			"$eq": true,
		},
	}
	cur, err := collection.Find(mango.MgCtx, filter)
	if err != nil {
		logging.Fatal(err)
	}
	var results []bson.M
	err = cur.All(mango.MgCtx, &results)
	if err != nil {
		logging.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func TestFindMany(t *testing.T) {
	var resultMe []ws.Trainer
	var resultYou []ws.Trainer
	// 获取collection
	sendIdCollection := mango.MgClient.Database("gin-IM").Collection("1-->2")
	idCollection := mango.MgClient.Database("gin-IM").Collection("2-->1")
	// SetSort 设置排序字段（1表示升序；-1表示降序）
	op1 := options.Find().SetSort(bson.D{{"start_time", -1}})
	op2 := options.Find().SetLimit(10)
	sendIdcur, err := sendIdCollection.Find(mango.MgCtx, bson.M{}, op1, op2)
	idcur, err := idCollection.Find(mango.MgCtx, bson.M{}, op1, op2)
	err = sendIdcur.All(mango.MgCtx, &resultYou) // sendId 对面发过来的
	err = idcur.All(mango.MgCtx, &resultMe)      // id 发给对面的
	if err != nil {
		logging.Info(err)
	}
	results, err := appendAndSort(resultMe, resultYou)
	if err != nil {
		logging.Info(err)
	}
	fmt.Println(results)
	fmt.Println("--------------------------------------")
	for _, result := range results {
		fmt.Println(result)
	}
}

// 对两个ws.Trainer类型的切片进行简单的排序和整合
func appendAndSort(resultMe []ws.Trainer, resultYou []ws.Trainer) ([]ws.Result, error) {
	var results []ws.Result
	for _, v := range resultMe {
		sendSort := SendSortMsg{ // 构造返回的msg
			Content:  v.Content,
			Read:     v.Read,
			CreateAt: v.StartTime,
		}
		result := ws.Result{ // 构造返回所有的内容,包括传送者
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

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     bool   `json:"read"`
	CreateAt int64  `json:"create_at"`
}
