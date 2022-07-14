package main

import (
	"fmt"
	"gin-IM/conf"
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
}
