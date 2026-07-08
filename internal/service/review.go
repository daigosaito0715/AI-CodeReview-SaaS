package service

import (
	"context"
	"fmt"
	"strings"

	"ai-codereview-saas/internal/model"
)

const reviewVersion = "dummy-v1"

type Reviewer interface {
	Review(ctx context.Context, req model.ReviewRequest) (model.ReviewResponse, error)
}

type DummyReviewService struct{}

func NewDummyReviewService() *DummyReviewService {
	return &DummyReviewService{}
}

func (s *DummyReviewService) Review(ctx context.Context, req model.ReviewRequest) (model.ReviewResponse, error) {
	issues := []model.ReviewIssue{
		{
			Severity:    "low",
			Title:       "println の使用",
			Description: "本格的なアプリケーションでは fmt.Println やログライブラリの使用を検討してください。",
			Location:    "コード全体",
			Reason:      "出力方法を標準的なAPIに寄せることで、後からログ出力やテストに置き換えやすくなります。",
		},
	}

	suggestions := []string{
		"fmt.Println または用途に合ったログ出力へ置き換える",
		"今後処理が増える場合に備えて、入出力処理とビジネスロジックを分ける",
	}

	return model.ReviewResponse{
		Language:      req.Language,
		ReviewVersion: reviewVersion,
		Score:         85,
		Summary:       "全体的にシンプルで読みやすいコードです。現時点では大きな問題はありませんが、出力方法と将来の拡張性に改善余地があります。",
		Strengths: []string{
			"コード量が少なく、処理の流れを追いやすい",
			"関数の役割が明確で、最小構成として理解しやすい",
			"不要な複雑さが少なく、改善方針を立てやすい",
		},
		Issues:      issues,
		Suggestions: suggestions,
		CodexPrompt: buildCodexPrompt(req, issues, suggestions),
	}, nil
}

func buildCodexPrompt(req model.ReviewRequest, issues []model.ReviewIssue, suggestions []string) string {
	var b strings.Builder

	fmt.Fprintf(&b, `以下のレビュー結果をもとに、%s のコードを改善してください。

目的:
- 既存の挙動を変えずに、読みやすさ・保守性・実運用時の扱いやすさを改善する
- レビュー指摘を反映し、必要以上に大きな設計変更は避ける

レビュー指摘:
`, req.Language)

	for _, issue := range issues {
		fmt.Fprintf(&b, "- [%s] %s: %s\n", issue.Severity, issue.Title, issue.Description)
		if issue.Location != "" {
			fmt.Fprintf(&b, "  該当箇所: %s\n", issue.Location)
		}
		if issue.Reason != "" {
			fmt.Fprintf(&b, "  改善理由: %s\n", issue.Reason)
		}
	}

	b.WriteString("\n改善案:\n")
	for _, suggestion := range suggestions {
		fmt.Fprintf(&b, "- %s\n", suggestion)
	}

	b.WriteString(`
作業条件:
- 既存の動作を維持する
- 変更理由を簡潔に説明する
- 可能であれば小さな差分で修正する
- 不要な依存関係を追加しない
`)

	return b.String()
}
