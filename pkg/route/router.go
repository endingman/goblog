package route

import (
	"github.com/gorilla/mux"
	"goblog/pkg/logger"
)

var Router *mux.Router

// Name2URL  通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := Router.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}
