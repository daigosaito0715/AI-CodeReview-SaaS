package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"ai-codereview-saas/internal/model"
	"ai-codereview-saas/internal/service"
)

const maxCodeLength = 100_000

type ReviewHandler struct {
	reviewer service.Reviewer
}

func NewReviewHandler(reviewer service.Reviewer) *ReviewHandler {
	return &ReviewHandler{reviewer: reviewer}
}

func (h *ReviewHandler) Review(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "許可されていないHTTPメソッドです。")
		return
	}

	var req model.ReviewRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "リクエストボディが不正なJSONです。")
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_json", "JSONは1つのオブジェクトだけを送信してください。")
		return
	}

	if err := validateReviewRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	result, err := h.reviewer.Review(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "review_failed", "レビュー処理に失敗しました。")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func validateReviewRequest(req model.ReviewRequest) error {
	if strings.TrimSpace(req.Language) == "" {
		return errors.New("language は必須です。")
	}
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("code は必須です。")
	}
	if len(req.Code) > maxCodeLength {
		return errors.New("code が長すぎます。")
	}
	return nil
}
