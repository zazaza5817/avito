package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
)

// userRepository реализует UserRepository
type userRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, user.UserID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update обновляет существующего пользователя
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, team_name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $4
	`
	result, err := r.db.ExecContext(ctx, query, user.Username, user.TeamName, user.IsActive, user.UserID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Get возвращает пользователя по ID
func (r *userRepository) Get(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByTeam возвращает всех пользователей команды
func (r *userRepository) GetByTeam(ctx context.Context, teamName string) ([]*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by team: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

// SetActive устанавливает статус активности пользователя
func (r *userRepository) SetActive(ctx context.Context, userID string, isActive bool) error {
	query := `
		UPDATE users
		SET is_active = $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, isActive, userID)
	if err != nil {
		return fmt.Errorf("failed to set user active status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetActiveTeammates возвращает активных участников команды, исключая указанного пользователя
func (r *userRepository) GetActiveTeammates(ctx context.Context, teamName string, excludeUserID string) ([]*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1 AND user_id != $2 AND is_active = true
		ORDER BY username
	`

	rows, err := r.db.QueryContext(ctx, query, teamName, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active teammates: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}
