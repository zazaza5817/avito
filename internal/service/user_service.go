package service

import (
	"context"
	"fmt"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/repository"
)

// UserService определяет интерфейс для работы с пользователями
type UserService interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserReviews(ctx context.Context, userID string) ([]*models.PullRequestShort, error)
}

// userService реализует UserService
type userService struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
}

// NewUserService создает новый сервис для работы с пользователями
func NewUserService(
	userRepo repository.UserRepository,
	prRepo repository.PullRequestRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

// SetUserActive устанавливает статус активности пользователя
func (s *userService) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	if err := s.userRepo.SetActive(ctx, userID, isActive); err != nil {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetUserReviews возвращает PR, где пользователь назначен ревьювером
func (s *userService) GetUserReviews(ctx context.Context, userID string) ([]*models.PullRequestShort, error) {
	// Проверяем существование пользователя
	_, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	prs, err := s.prRepo.GetByReviewer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}

	return prs, nil
}
