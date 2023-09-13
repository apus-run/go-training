package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/apus-run/sea-kit/slogx"

	"gin-with-render/internal/web"
	"gin-with-render/internal/web/handler"
)

func main() {
	// logger := slogx.NewLogger(slogx.WithEncoding("json"), slogx.WithFilename("logs.log"))
	logger := slogx.NewLogger(slogx.WithFilename(""), slogx.WithEncoding("console"))
	//logger.Debug("This is a debug message", slog.Any("key", "value"))
	//logger.Info("This is a info message")
	//logger.Warn("This is a warn message")
	//logger.Error("This is a error message")
	slog.SetDefault(logger)

	log.Printf("Starting up with pid %d...", os.Getpid())

	userHandler := handler.NewUserHandler(logger)
	router := web.Router(userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	gl := logger.WithGroup("server")
	gl.Info("后端服务启动了", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
