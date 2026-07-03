package service

import (
	"context"

	"ai-codereview-saas/internal/model"
)

type Reviewer interface {
	Review(ctx context.Context, req model.ReviewRequest) (model.ReviewResponse, error)
}

type DummyReviewService struct{}

func NewDummyReviewService() *DummyReviewService {
	return &DummyReviewService{}
}

func (s *DummyReviewService) Review(ctx context.Context, req model.ReviewRequest) (model.ReviewResponse, error) {
	return model.ReviewResponse{
		Score:   85,
		Summary: "全体的にシンプルで読みやすいコードです。",
		Issues: []model.ReviewIssue{
			{
				Severity:    "low",
				Title:       "println の使用",
				Description: "本格的なアプリケーションでは fmt.Println やログライブラリの使用を検討してください。",
			},
		},
		CodexPrompt: "以下のレビュー内容をもとに、既存の挙動を変えずにコードを改善してください。",
	}, nil
}
