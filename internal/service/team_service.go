package service

import (
	"context"
	"fmt"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/repository"
)

// TeamService определяет интерфейс для работы с командами
type TeamService interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
}

// teamService реализует TeamService
type teamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

// NewTeamService создает новый сервис для работы с командами
func NewTeamService(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// CreateTeam создает команду с участниками
func (s *teamService) CreateTeam(ctx context.Context, team *models.Team) error {
	// Проверяем, существует ли команда
	exists, err := s.teamRepo.Exists(ctx, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		return ErrTeamExists
	}

	// Создаем команду
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	// Создаем или обновляем участников
	for _, member := range team.Members {
		user := &models.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: team.TeamName,
			IsActive: member.IsActive,
		}

		// Пытаемся получить существующего пользователя
		existingUser, err := s.userRepo.Get(ctx, member.UserID)
		if err != nil {
			// Если пользователь не найден, создаем нового
			if err := s.userRepo.Create(ctx, user); err != nil {
				return fmt.Errorf("failed to create user %s: %w", member.UserID, err)
			}
		} else {
			// Если пользователь существует, обновляем его
			existingUser.Username = user.Username
			existingUser.TeamName = user.TeamName
			existingUser.IsActive = user.IsActive
			if err := s.userRepo.Update(ctx, existingUser); err != nil {
				return fmt.Errorf("failed to update user %s: %w", member.UserID, err)
			}
		}
	}

	return nil
}

// GetTeam возвращает команду с участниками
func (s *teamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	team, err := s.teamRepo.Get(ctx, teamName)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
