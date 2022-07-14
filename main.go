package main

import (
	config "gin-IM/conf"
	"gin-IM/db/mango"
	"gin-IM/db/mysql"
	"gin-IM/db/redis"
	"gin-IM/route"

	log "github.com/sirupsen/logrus"
)

func main() {
	r := route.NewRoute()
	err := r.Run(config.HttpPort)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	config.Init("./conf/config.ini")
	mysql.Init()
	redis.Init()
	mango.Init()
}
