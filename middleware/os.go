package middleware

import (
	"context"
	ua "github.com/mileusna/useragent"
	"net/http"
)

func Os(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//OSを解析
		type UserAgentKey struct{}
		const OS = "OS"
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), OS, ua.Parse(r.UserAgent()).OS)))
	}
	return http.HandlerFunc(fn)
}
