package mysql

import (
	config "gin-IM/conf"
	"gin-IM/model"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var MysqlDB *gorm.DB // 全局MysqlDB

func Init() {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	var builder strings.Builder
	s := []string{config.MysqlUser, ":", config.MysqlPassword, "@tcp(", config.MysqlHost, ":", config.MysqlPort, ")/", config.MysqlName, "?charset=utf8mb4&parseTime=True&loc=Local"}
	for _, str := range s {
		builder.WriteString(str)
	}
	dsn := builder.String()
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Infoln(err)
	}
	db.LogMode(config.MysqlIsLog)                           // 开启 Logger, 以展示详细的日志
	db.SingularTable(config.MysqlIsSingularTable)           // 如果设置禁用表名复数形式属性为 true，`User` 的表名将是 `user`(因为gorm默认表名是复数)
	db.DB().SetMaxIdleConns(config.MysqlMaxIdleConns)       // 设置空闲连接池中的最大连接数
	db.DB().SetMaxOpenConns(config.MysqlMaxOpenConns)       // 设置数据库连接最大打开数。
	db.DB().SetConnMaxLifetime(config.MysqlConnMaxLifetime) // 设置可重用连接的最长时间
	// 自动迁移
	db.AutoMigrate(&model.User{})
	MysqlDB = db
}
