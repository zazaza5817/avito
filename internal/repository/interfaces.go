// Package repository определяет интерфейсы для слоя доступа к данным.
package repository

import (
	"context"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
)

// TeamRepository определяет интерфейс для работы с командами
type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	Get(ctx context.Context, teamName string) (*models.Team, error)
	Exists(ctx context.Context, teamName string) (bool, error)
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Get(ctx context.Context, userID string) (*models.User, error)
	GetByTeam(ctx context.Context, teamName string) ([]*models.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) error
	GetActiveTeammates(ctx context.Context, teamName string, excludeUserID string) ([]*models.User, error)
}

// PullRequestRepository определяет интерфейс для работы с Pull Request
type PullRequestRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	Get(ctx context.Context, prID string) (*models.PullRequest, error)
	Update(ctx context.Context, pr *models.PullRequest) error
	Exists(ctx context.Context, prID string) (bool, error)
	GetByReviewer(ctx context.Context, reviewerID string) ([]*models.PullRequestShort, error)
	AssignReviewer(ctx context.Context, prID, reviewerID string) error
	RemoveReviewer(ctx context.Context, prID, reviewerID string) error
	GetReviewers(ctx context.Context, prID string) ([]string, error)
	IsReviewerAssigned(ctx context.Context, prID, reviewerID string) (bool, error)
}
