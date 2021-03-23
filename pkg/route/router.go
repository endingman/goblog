package route

import (
	"github.com/gorilla/mux"
	"goblog/pkg/logger"
	"net/http"
)

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter()
}

// Name2URL  通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
