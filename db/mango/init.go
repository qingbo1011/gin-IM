package mango

import (
	"context"
	config "gin-IM/conf"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MgCtx = context.Background() // 全局MangoDB ctx
var MgClient *mongo.Client       // 全局MangoDB Client

func Init() {
	// mongoClient 连接客户端参数
	clientOpts := options.Client().
		SetAuth(options.Credential{
			AuthMechanism: config.MangoAuthMechanism,
			//AuthSource:              "anquan",	// 用于身份验证的数据库的名称
			Username: config.MangoUser,
			Password: config.MangoPassword,
		}).
		SetConnectTimeout(config.MangoConnectTimeout).
		SetHosts(config.MangoHosts).
		SetMaxPoolSize(config.MangoMaxPoolSize).
		SetMinPoolSize(config.MangoMinPoolSize).
		SetReadPreference(readpref.Primary()). // 默认值是readpref.Primary()（https://www.mongodb.com/docs/manual/core/read-preference/#read-preference）
		SetReplicaSet("")                      // SetReplicaSet指定集群的副本集名称。（默认为空）

	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Infoln(err)
	}
	MgClient = client
}
