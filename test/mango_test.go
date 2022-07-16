package test

import (
	"fmt"
	config "gin-IM/conf"
	"gin-IM/db/mango"
	"testing"

	log "github.com/sirupsen/logrus"
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
		log.Fatal(err)
	}
	var results []bson.M
	err = findCursor.All(mango.MgCtx, &results)
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
