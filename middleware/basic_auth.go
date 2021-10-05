package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		os.Setenv("BASIC_AUTH_USER_ID", "id")
		os.Setenv("BASIC_AUTH_PASSWORD", "pass")
		authId := os.Getenv("BASIC_AUTH_USER_ID")
		authPass := os.Getenv("BASIC_AUTH_PASSWORD")
		userID, password, ok := r.BasicAuth()
		if ok == false {
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
