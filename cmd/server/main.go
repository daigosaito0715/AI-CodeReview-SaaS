package main

import (
	"log"
	"net/http"
	"os"

	"ai-codereview-saas/internal/handler"
	"ai-codereview-saas/internal/service"
)

func main() {
	reviewService := service.NewDummyReviewService()
	reviewHandler := handler.NewReviewHandler(reviewService)

	router := handler.NewRouter(reviewHandler)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
