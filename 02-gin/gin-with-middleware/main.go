package main

import (
	"log"
	"net/http"
	"time"

	"gin-with-middleware/router"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Router(),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server listen at %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
