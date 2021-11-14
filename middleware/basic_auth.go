package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	authId := os.Getenv("BASIC_AUTH_USER_ID")
	authPass := os.Getenv("BASIC_AUTH_PASSWORD")
	fn := func(w http.ResponseWriter, r *http.Request) {
		userID, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if userID == authId && password == authPass {
			h.ServeHTTP(w, r)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	return http.HandlerFunc(fn)
}
