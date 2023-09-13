package main

import (
	"log"
	"net/http"
	"time"

	"gin-with-seesion/internal/web"
	"gin-with-seesion/internal/web/handler"
)

func main() {
	userHandler := handler.NewUserHandler()
	router := web.Router(userHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server listen at %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
