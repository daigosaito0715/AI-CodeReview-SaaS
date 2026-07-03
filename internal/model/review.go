package model

type ReviewRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ReviewResponse struct {
	Score       int           `json:"score"`
	Summary     string        `json:"summary"`
	Issues      []ReviewIssue `json:"issues"`
	CodexPrompt string        `json:"codex_prompt"`
}

type ReviewIssue struct {
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
