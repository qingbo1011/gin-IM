package test

import (
	"fmt"
	config "gin-IM/conf"
	"gin-IM/db/redis"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	config.Init("../conf/config.ini")
	redis.Init()
}

func TestRedis(t *testing.T) {
	result, err := redis.Rdb.Get(redis.Rctx, "hello").Result()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}
