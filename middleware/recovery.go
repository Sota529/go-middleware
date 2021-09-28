package middleware

import (
	"fmt"
	"net/http"
)

type RecoveryHandler struct{}

func NewRecoveryHandler() *RecoveryHandler {
	return &RecoveryHandler{}
}

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recoverd!")
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (R *RecoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic!")
}
