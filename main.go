package main

import (
	"fmt"
	config "gin-IM/conf"
	"gin-IM/db/mysql"
)

func main() {
	fmt.Println(config.HttpPort)
	fmt.Println(config.MysqlHost)
	fmt.Println(config.RedisDbName)
	fmt.Println(config.MangoHosts)
	fmt.Println(config.MangoConnectTimeout)

}

func init() {
	config.Init("./conf/config.ini")
	mysql.Init()
}
