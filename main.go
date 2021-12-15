package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/middleware"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

var wg sync.WaitGroup

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// set http handlers
	mux := http.NewServeMux()
	healthz := handler.NewHealthzHandler()

	svc := service.NewTODOService(todoDB)
	todos := handler.NewTODOHandler(svc)

	// TODO: ここから実装を行う
	mux.Handle("/healthz", healthz)
	mux.Handle("/todos", todos)
	mux.Handle("/do-panic", middleware.Os(middleware.ProcessTime(http.HandlerFunc(HeaveFunc))))
	mux.Handle("/auth", middleware.BasicAuth(http.HandlerFunc(HeaveFunc)))
	srv := &http.Server{
		Addr:    defaultPort,
		Handler: middleware.Recovery(mux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalln("Server closed with error:", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Println("Failed to gracefully shutdown:", err)
	}
	log.Println("Server shutdown")
	return nil
}

func HeaveFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("heavy process starts")
	time.Sleep(3 * time.Second)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("done!\n"))
}
