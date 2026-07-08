package model

type ReviewRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ReviewResponse struct {
	Language      string        `json:"language"`
	ReviewVersion string        `json:"review_version"`
	Score         int           `json:"score"`
	Summary       string        `json:"summary"`
	Strengths     []string      `json:"strengths"`
	Issues        []ReviewIssue `json:"issues"`
	Suggestions   []string      `json:"suggestions"`
	CodexPrompt   string        `json:"codex_prompt"`
}

type ReviewIssue struct {
	Severity    string `json:"severity"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Reason      string `json:"reason"`
}
