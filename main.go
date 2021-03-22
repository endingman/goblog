package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

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
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "创建新的文章")
}

func main() {
	router := mux.NewRouter()//mux路由，gorilla/mux 因实现了 net/http 包的 http.Handler 接口，故兼容 http.ServeMux

	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
    router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

    router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
    router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
    router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

    // 自定义 404 页面
    router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

    // 通过命名路由获取 URL 示例
    homeURL, _ := router.Get("home").URL()
    fmt.Println("homeURL: ", homeURL)
    articleURL, _ := router.Get("articles.show").URL("id", "23")
    fmt.Println("articleURL: ", articleURL)

    http.ListenAndServe(":3000", router)

	//http.ListenAndServe 用以监听本地 3000 端口以提供服务，标准的 HTTP 端口是 80 端口
	http.ListenAndServe(":3000", router)
}
