package redis

import (
	"context"
	config "gin-IM/conf"

	"github.com/go-redis/redis/v8"
)

var Rctx = context.Background() // 全局Redis ctx
var Rdb *redis.Client           // 全局Redis DB

func Init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,     // redis地址
		Password: config.RedisPassword, // redis密码，没有则留空
		DB:       config.RedisDbName,   // 默认数据库（不指定默认是0）
	})
}
