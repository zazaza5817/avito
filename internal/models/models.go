package models

import "time"

// User представляет пользователя системы
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// Team представляет команду с участниками
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// TeamMember представляет участника команды
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// PullRequestStatus представляет статус Pull Request
type PullRequestStatus string

// Константы статусов Pull Request
const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)

// PullRequest представляет Pull Request
type PullRequest struct {
	PullRequestID     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	AuthorID          string            `json:"author_id"`
	Status            PullRequestStatus `json:"status"`
	AssignedReviewers []string          `json:"assigned_reviewers"`
	CreatedAt         *time.Time        `json:"createdAt,omitempty"`
	MergedAt          *time.Time        `json:"mergedAt,omitempty"`
}

// PullRequestShort представляет краткую информацию о Pull Request
type PullRequestShort struct {
	PullRequestID   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	AuthorID        string            `json:"author_id"`
	Status          PullRequestStatus `json:"status"`
}
