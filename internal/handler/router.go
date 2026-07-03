package handler

import "net/http"

func NewRouter(reviewHandler *ReviewHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/review", reviewHandler.Review)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "指定されたエンドポイントは存在しません。")
	})
	return mux
}
