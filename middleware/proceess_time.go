package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/TechBowl-japan/go-stations/model"
	"log"
	"net/http"
	"time"
)

func ProcessTime(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		//処理時間測定
		end := time.Now()
		Log, err := json.Marshal(model.Access{
			Timestamp: startTime,
			Latency:   (end.Sub(startTime)).Milliseconds(),
			Path:      r.RequestURI,
			OS:        r.Context().Value("OS").(string),
		})
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(Log))
	}
	return http.HandlerFunc(fn)
}
