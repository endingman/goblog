package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"net/http"
	"net/url"
	/**
	因为引入的是驱动，操作数据库时我们使用的是 sql 库里的方法，而不会具体使用到 go-sql-driver/mysql 包里的方法，
	当有未使用的包被引入时，Go 编译器会停止编译。
	为了让编译器能正常运行，需要使用 匿名导入 来加载。
	*/
	_ "github.com/go-sql-driver/mysql"
)

//包级别的变量声明时不能使用 := 语法，修改为带关键词 var 的变量声明
var router = mux.NewRouter().StrictSlash(true)
var db *sql.DB

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

type Article struct {
	Title, Body string
	ID          int64
}

func main() {
	//mux路由，gorilla/mux 因实现了 net/http 包的 http.Handler 接口，故兼容 http.ServeMux
	//访问以下两个链接：
	//localhost:3000/about
	//localhost:3000/about/
	//可以看到有 / 的链接会报 404 错误：
	//Gorilla Mux 提供了一个 StrictSlash(value bool) 函数处理`/`问题
	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()

	router := bootstrap.SetupRoute()

	// 通过命名路由获取 URL 示例
	//homeURL, _ := router.Get("home").URL()
	//fmt.Println("homeURL: ", homeURL)
	//articleURL, _ := router.Get("articles.show").URL("id", "23")
	//fmt.Println("articleURL: ", articleURL)

	//http.ListenAndServe 用以监听本地 3000 端口以提供服务，标准的 HTTP 端口是 80 端口
	http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
}
