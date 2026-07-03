package handler

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "許可されていないHTTPメソッドです。")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
