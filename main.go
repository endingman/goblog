package main

import (
	"fmt"
	"net/http"
)

/**
http.HandleFunc 用以指定处理 HTTP 请求的函数，
此函数允许我们只写一个 handler（在此例子中 handlerFunc，可任意命名），
请求会通过参数传递进来，使用者只需与 http.Request 和 http.ResponseWriter 两个对象交互即可。

http.Request 是用户的请求信息，一般用 r 作为简写。
http.ResponseWriter 是返回用户的响应，一般用 w 作为简写。
*/
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello , 欢迎来到 goblog！</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
			"<p>如有疑惑，请联系我们。</p>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func main() {
	router := http.NewServeMux()

	//方法过滤
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, "访问文章列表")
		case "POST":
			fmt.Fprint(w, "创建新的文章")
		}
	})

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/about", aboutHandler)
	//http.ListenAndServe 用以监听本地 3000 端口以提供服务，标准的 HTTP 端口是 80 端口
	http.ListenAndServe(":3000", router)
}
