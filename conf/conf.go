package config

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	HttpPort string

	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassWord string
	MysqlName     string

	RedisAddr     string
	RedisPassWord string
	RedisDbName   int

	MangoAuthMechanism  string
	MangoUser           string
	MangoPassword       string
	MangoHosts          string
	MangoConnectTimeout time.Duration
	MangoMaxPoolSize    int
	MangoMinPoolSize    int
)

func Init(path string) {
	file, err := ini.Load(path)
	if err != nil {
		log.Fatalln("Fail to parse 'conf/app.ini': ", err)
	}

	loadService(file)
	loadMysql(file)
	loadRedis(file)
	loadMongo(file)
}

func loadService(file *ini.File) {
	HttpPort = file.Section("service").Key("HttpPort").MustString(":8080")
}

func loadMysql(file *ini.File) {
	section, err := file.GetSection("mysql")
	if err != nil {
		log.Fatalln(err)
	}
	MysqlHost = section.Key("MysqlHost").String()
	MysqlPort = section.Key("MysqlPort").String()
	MysqlUser = section.Key("MysqlUser").String()
	MysqlPassWord = section.Key("MysqlPassWord").String()
	MysqlName = section.Key("MysqlName").String()
}

func loadRedis(file *ini.File) {
	section, err := file.GetSection("redis")
	if err != nil {
		log.Fatalln(err)
	}
	RedisAddr = section.Key("RedisAddr").String()
	RedisPassWord = section.Key("RedisPassWord").String()
	RedisDbName = section.Key("RedisDbName").MustInt(1) // MustInt，defaultVal为1
}

func loadMongo(file *ini.File) {
	section, err := file.GetSection("mongo")
	if err != nil {
		log.Fatalln(err)
	}
	MangoAuthMechanism = section.Key("MangoAuthMechanism").String()
	MangoUser = section.Key("MangoUser").String()
	MangoPassword = section.Key("MangoPassword").String()
	MangoHosts = section.Key("MangoHosts").String()
	MangoConnectTimeout = time.Duration(section.Key("MangoConnectTimeout").MustInt(10)) * time.Second
	MangoMaxPoolSize = section.Key("MangoMaxPoolSize").MustInt(20)
	MangoMinPoolSize = section.Key("MangoMinPoolSize").MustInt(5)
}
