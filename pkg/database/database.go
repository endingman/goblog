package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"goblog/pkg/logger"
	"time"
)

var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	initDB()
	createTables()
}

func initDB() {
	var err error

	/**
	DSN 全称为 Data Source Name，表示 数据源信息，用于定义如何连接数据库。
	不同数据库的 DSN 格式是不同的，这取决于数据库驱动的实现，
	下面是 go-sql-driver/sql 的 DSN 格式，如下所示：
	//[用户名[:密码]@][协议(数据库服务器地址)]]/数据库名称?参数列表
	[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	*/

	config := mysql.Config{
		User:                 "homestead",
		Passwd:               "secret",
		Addr:                 "127.0.0.1:33060",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	// 准备数据库连接池
	//config.FormatDSN() homestead:secret@tcp(127.0.0.1:33060)/goblog?checkConnLiveness=false&maxAllowedPacket=0
	//func Open(driverName, dataSourceName string) (*sql.DB, error)
	DB, err = sql.Open("mysql", config.FormatDSN())

	logger.LogError(err)

	// 设置最大连接数
	//实验表明，在高并发的情况下，将值设为大于 10，可以获得比设置为 1 接近六倍的性能提升。
	// 而设置为 10 跟设置为 0（也就是无限制），在高并发的情况下，性能差距不明显。
	//需要考虑的是不要超出数据库系统设置的最大连接数。
	//show variables like 'max_connections';
	DB.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	/**
	设置连接池最大空闲数据库连接数，<= 0 表示不设置空闲连接数，默认为 2。

	实验表明，在高并发的情况下，将值设为大于 0，可以获得比设置为 0 超过 20 倍的性能提升。

	这是因为设置为 0 的情况下，每一个 SQL 连接执行任务以后就销毁掉了，执行新任务时又需要重新建立连接。

	很明显，重新建立连接是很消耗资源的一个动作。
	*/
	DB.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
   ); `
	/**
	Exec() 来执行创建数据库表结构的语句。
	一般使用 sql.DB 中的 Exec() 来执行没有返回结果集的 SQL 语句
	*/
	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
