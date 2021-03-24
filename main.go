package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"net/url"
	"strings"
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

	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

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

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)
	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 SQL 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}

func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + types.Int64ToString(a.ID))

	if err != nil {
		return 0, err
	}

	// √ 删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)
	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示表单
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
	}
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误

		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {

			// 4.2 表单验证通过，更新数据
			/**
			Exec () 方法
			执行数据更新的是 Exec() 方法，此方法与我们之前学习 Prepare 方法时搭配使用 stmt.Exec() 不一样，
			stmt.Exec() 是 sql.Stmt 的方法，而这里的 Exec() 是 sql.DB 提供的方法。

			一般情况下，我们使用此方法来处理 CREATE、UPDATE、DELETE 类型的 SQL。

			与 createTables() 方法中使用的 Exec() 一致
			*/
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			// √ 更新成功，跳转到文章详情页
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {

			// 4.3 表单验证不通过，显示理由
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)

			tmpl.Execute(w, data)
		}
	}
}

func getArticleByID(id string) (Article, error) {
	// QueryRow() 是可变参数的方法，参数可以为一个或者多个。
	// 参数只有一个的情况下，我们称之为纯文本模式，多个参数的情况下称之为 Prepare 模式。
	// QueryRow() 封装了 Prepare 方法的调用
	/**
	QueryRow() 会返回一个 sql.Row struct，紧接着我们使用链式调用的方式调用了 sql.Row.Scan() 方法：
	db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	Scan() 将查询结果赋值到我们的 article struct 中，传参应与数据表字段的顺序保持一致。
	*/
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
