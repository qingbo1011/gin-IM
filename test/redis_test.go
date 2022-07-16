package test

import (
	"fmt"
	config "gin-IM/conf"
	"gin-IM/db/redis"
	"testing"

	logging "github.com/sirupsen/logrus"
)

func init() {
	config.Init("../conf/config.ini")
	redis.Init()
}

func TestRedis(t *testing.T) {
	result, err := redis.Rdb.Get(redis.RCtx, "hello").Result()
	if err != nil {
		logging.Fatalln(err)
	}
	fmt.Println(result)
}

func TestCreateHashRedis(t *testing.T) {
	m := make(map[string]any, 0)
	m["User-Agent2"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzNzAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTczMDF9.ipiIDgAdTwrv8EX45y0UD6wy0fOOdzhIDysyB8kJais"
	result, err := redis.Rdb.HSet(redis.RCtx, "uid", m).Result()
	if err != nil {
		logging.Info(err)
	}
	fmt.Println(result)
}

func TestSelectRedis(t *testing.T) {
	val, err := redis.Rdb.HGet(redis.RCtx, "1", "PostmanRuntime/7.26.8").Result()
	if err != nil {
		logging.Info(err)
	}
	fmt.Println(val)
}
