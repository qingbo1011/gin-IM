package main

import (
	"fmt"
	config "gin-IM/conf"
	"gin-IM/db/mango"
	"gin-IM/db/mysql"
	"gin-IM/db/redis"
)

func main() {
	fmt.Println(config.MangoHosts)
	fmt.Println(config.HttpPort)
	fmt.Printf("%T", config.MangoHosts)
}

func init() {
	config.Init("./conf/config.ini")
	mysql.Init()
	redis.Init()
	mango.Init()
}
