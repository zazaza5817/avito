package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
)

// teamRepository реализует TeamRepository
type teamRepository struct {
	db *sql.DB
}

// NewTeamRepository создает новый репозиторий команд
func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

// Create создает новую команду
func (r *teamRepository) Create(ctx context.Context, team *models.Team) error {
	query := `INSERT INTO teams (team_name) VALUES ($1)`
	_, err := r.db.ExecContext(ctx, query, team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

// Get возвращает команду с участниками
func (r *teamRepository) Get(ctx context.Context, teamName string) (*models.Team, error) {
	exists, err := r.Exists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("team not found")
	}

	query := `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var members []models.TeamMember
	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan team member: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &models.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

// Exists проверяет существование команды
func (r *teamRepository) Exists(ctx context.Context, teamName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, teamName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}
	return exists, nil
}
