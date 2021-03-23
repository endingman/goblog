package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
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

/**
http.HandleFunc 用以指定处理 HTTP 请求的函数，
此函数允许我们只写一个 handler（在此例子中 handlerFunc，可任意命名），
请求会通过参数传递进来，使用者只需与 http.Request 和 http.ResponseWriter 两个对象交互即可。

http.Request 是用户的请求信息，一般用 r 作为简写。
http.ResponseWriter 是返回用户的响应，一般用 w 作为简写。
*/
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	//	获取URL参数
	vars := mux.Vars(r)
	id := vars["id"]

	//读取对应文章数据
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	// QueryRow() 是可变参数的方法，参数可以为一个或者多个。
	// 参数只有一个的情况下，我们称之为纯文本模式，多个参数的情况下称之为 Prepare 模式。
	// QueryRow() 封装了 Prepare 方法的调用
	/**
	QueryRow() 会返回一个 sql.Row struct，紧接着我们使用链式调用的方式调用了 sql.Row.Scan() 方法：
	db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	Scan() 将查询结果赋值到我们的 article struct 中，传参应与数据表字段的顺序保持一致。
	*/
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	//	如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//	数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示文章
		tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w, article)
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	//Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。
	//PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。

	/**
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, "请提供正确的参数")
		return
	}

	title := r.PostForm.Get("title")


	fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	fmt.Fprintf(w, "title 的值为: %v", title)

	//如不想获取所有的请求内容，而是逐个获取的话，这也是比较常见的操作，
	// 无需使用 r.ParseForm() 可直接使用 r.FormValue() 和 r.PostFormValue() 方法
	fmt.Fprintf(w, "r.Form 中 title 的值为: %v <br>", r.FormValue("title"))
	fmt.Fprintf(w, "r.PostForm 中 title 的值为: %v <br>", r.PostFormValue("title"))
	fmt.Fprintf(w, "r.Form 中 test 的值为: %v <br>", r.FormValue("test"))
	fmt.Fprintf(w, "r.PostForm 中 test 的值为: %v <br>", r.PostFormValue("test"))
	*/

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := make(map[string]string)

	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
		/**
		Go 语言的内建函数 len ()，可以用来获取切片、字符串、通道（channel）等的长度。

		这里的差异是由于 Go 语言的字符串都以 UTF-8 格式保存，每个中文占用 3 个字节，因此使用 len () 获得两个中文文字对应的 6 个字节。

		如果希望按习惯上的字符个数来计算，就需要使用 Go 语言中 utf8 包提供的 RuneCountInString () 函数来计数
		*/
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	// 检查是否有错误
	if len(errors) == 0 {
		fmt.Fprint(w, "验证通过!<br>")
		fmt.Fprintf(w, "title 的值为: %v <br>", title)
		fmt.Fprintf(w, "title 的长度为: %v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body 的值为: %v <br>", body)
		fmt.Fprintf(w, "body 的长度为: %v <br>", utf8.RuneCountInString(body))

		lastInsertId, err := saveArticleToDB(title, body)
		if lastInsertId > 0 {
			// Go 标准库的 strconv 包。此包主要提供字符串和其他类型之间转换的函数。
			//类型转换在脚本类语言例如说 PHP 或者 JS 中不需要太重视，
			// 但在 Go 强类型语言中是一个很重要的概念。
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertId, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}

	} else {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

func saveArticleToDB(title string, body string) (int64, error) {
	/**
	多变量声明的方式与引入多个包使用 import(...) 同出一辙，
	都是 Go 语言为了让开发者少写代码而提供的简写方式。
	*/
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	//	获取一个prepare
	// 在数据库安全方面，Prepare 语句是防范 SQL 注入攻击有效且必备的手段。
	// 可以理解为将包含变量占位符 ? 的语句先告知 MySQL 服务器端。
	// Prepare 只会生产 stmt ，真正执行请求的需要调用 stmt.Exec()
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	//	例行错误检测
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	//执行请求
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	//插入成功的话返回自增ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, err
	}

	return 0, err
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		/**
		在 Go 中，一般 err 处理方式可以是给用户提示或记录到错误日志里，这种很多时候为 业务逻辑错误。
		当有重大错误，或者系统错误时，例如无法加载模板文件，就使用 panic() 。
		应用里需要有一套合理的错误机制，后面的开发中我们会详细讲解到。
		*/
		panic(err)
	}

	tmpl.Execute(w, data)
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2. 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			//strings 包提供的 TrimSuffix(s, suffix string) string 函数来移除 / 后缀
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	//mux路由，gorilla/mux 因实现了 net/http 包的 http.Handler 接口，故兼容 http.ServeMux
	//访问以下两个链接：
	//localhost:3000/about
	//localhost:3000/about/
	//可以看到有 / 的链接会报 404 错误：
	//Gorilla Mux 提供了一个 StrictSlash(value bool) 函数处理`/`问题
	initDB()
	createTables()

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	//homeURL, _ := router.Get("home").URL()
	//fmt.Println("homeURL: ", homeURL)
	//articleURL, _ := router.Get("articles.show").URL("id", "23")
	//fmt.Println("articleURL: ", articleURL)

	//http.ListenAndServe 用以监听本地 3000 端口以提供服务，标准的 HTTP 端口是 80 端口
	http.ListenAndServe(":3000", removeTrailingSlash(router))
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
	db, err = sql.Open("mysql", config.FormatDSN())

	checkError(err)

	// 设置最大连接数
	//实验表明，在高并发的情况下，将值设为大于 10，可以获得比设置为 1 接近六倍的性能提升。
	// 而设置为 10 跟设置为 0（也就是无限制），在高并发的情况下，性能差距不明显。
	//需要考虑的是不要超出数据库系统设置的最大连接数。
	//show variables like 'max_connections';
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	/**
	设置连接池最大空闲数据库连接数，<= 0 表示不设置空闲连接数，默认为 2。

	实验表明，在高并发的情况下，将值设为大于 0，可以获得比设置为 0 超过 20 倍的性能提升。

	这是因为设置为 0 的情况下，每一个 SQL 连接执行任务以后就销毁掉了，执行新任务时又需要重新建立连接。

	很明显，重新建立连接是很消耗资源的一个动作。
	*/
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	// 尝试连接，失败会报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}
