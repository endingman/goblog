package routes

import (
	"github.com/gorilla/mux"
	"goblog/app/http/controllers"
	"goblog/app/http/middlewares"
	"net/http"
)

// RegisterWebRoutes 注册网页相关路由
func RegisterWebRoutes(r *mux.Router) {
	//静态页面
	pc := new(controllers.PagesController)
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	//	articles
	ac := new(controllers.ArticlesController)
	r.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	// 静态资源
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 中间件：强制内容类型为 HTML
	//r.Use(middlewares.ForceHTML)

	// 用户认证
	auc := new(controllers.AuthController)
	r.HandleFunc("/auth/register", auc.Register).Methods("GET").Name("auth.register")
	r.HandleFunc("/auth/do-register", auc.DoRegister).Methods("POST").Name("auth.doregister")
	r.HandleFunc("/auth/login", auc.Login).Methods("GET").Name("auth.login")
	r.HandleFunc("/auth/dologin", auc.DoLogin).Methods("POST").Name("auth.dologin")
	r.HandleFunc("/auth/logout", auc.Logout).Methods("POST").Name("auth.logout")

	// 开始会话
	r.Use(middlewares.StartSession)
}
