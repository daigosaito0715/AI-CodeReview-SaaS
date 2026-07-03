package main

import (
	"log"
	"net/http"

	"ai-codereview-saas/internal/handler"
	"ai-codereview-saas/internal/service"
)

func main() {
	reviewService := service.NewDummyReviewService()
	reviewHandler := handler.NewReviewHandler(reviewService)

	router := handler.NewRouter(reviewHandler)

	addr := ":8080"
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
