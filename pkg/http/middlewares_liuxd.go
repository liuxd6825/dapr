package http

import (
	"github.com/liuxd6825/dapr/pkg/security/consts"
	"net/http"
)

// APIServiceInvoke 服务调用
func APIServiceInvoke(token string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if token == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Header.Get(consts.APITokenHeader)
			if v != token && !isRouteExcludedFromAPITokenAuth(r.Method, r.URL) {
				http.Error(w, "invalid api token", http.StatusUnauthorized)
				return
			}

			r.Header.Del(consts.APITokenHeader)
			next.ServeHTTP(w, r)
		})
	}
}
