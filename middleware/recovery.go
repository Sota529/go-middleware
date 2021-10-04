package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			//panicリカバリー
			if err := recover(); err != nil {
				fmt.Println("recoverd!")
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func HandlePanic(w http.ResponseWriter, r *http.Request) {
	//panic("panic!")
}
