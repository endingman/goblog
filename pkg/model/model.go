package model

import (
	"goblog/pkg/logger"

	"gorm.io/gorm"

	// GORM 的 MYSQL 数据库驱动导入
	"gorm.io/driver/mysql"

	gormlogger "gorm.io/gorm/logger"
)

// DB gorm.DB 对象
var DB *gorm.DB

// ConnectDB 初始化模型
func ConnectDB() *gorm.DB {

	var err error

	config := mysql.New(mysql.Config{
		DSN: "homestead:secret@tcp(127.0.0.1:33060)/goblog?charset=utf8&parseTime=True&loc=Local",
	})

	// 准备数据库连接池
	DB, err = gorm.Open(config, &gorm.Config{
		//开启gorm日志
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})

	logger.LogError(err)

	return DB
}
