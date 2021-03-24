package middlewares

import (
	"net/http"
	"strings"
)

func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			//strings 包提供的 TrimSuffix(s, suffix string) string 函数来移除 / 后缀
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
