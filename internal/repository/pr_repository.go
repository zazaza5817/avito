package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
)

// prRepository реализует PullRequestRepository
type prRepository struct {
	db *sql.DB
}

// NewPullRequestRepository создает новый репозиторий Pull Request
func NewPullRequestRepository(db *sql.DB) PullRequestRepository {
	return &prRepository{db: db}
}

// Create создает новый Pull Request
func (r *prRepository) Create(ctx context.Context, pr *models.PullRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Создаем PR
	query := `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	now := time.Now()
	_, err = tx.ExecContext(ctx, query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, now)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	// Назначаем ревьюверов
	for _, reviewerID := range pr.AssignedReviewers {
		err = r.assignReviewerTx(ctx, tx, pr.PullRequestID, reviewerID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	pr.CreatedAt = &now
	return nil
}

// Get возвращает Pull Request по ID
func (r *prRepository) Get(ctx context.Context, prID string) (*models.PullRequest, error) {
	query := `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var pr models.PullRequest
	var createdAt time.Time
	var mergedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, prID).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&pr.Status,
		&createdAt,
		&mergedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("pull request not found")
		}
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	pr.CreatedAt = &createdAt
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	// Получаем ревьюверов
	reviewers, err := r.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

// Update обновляет Pull Request
func (r *prRepository) Update(ctx context.Context, pr *models.PullRequest) error {
	query := `
		UPDATE pull_requests
		SET pull_request_name = $1, status = $2, merged_at = $3
		WHERE pull_request_id = $4
	`

	result, err := r.db.ExecContext(ctx, query, pr.PullRequestName, pr.Status, pr.MergedAt, pr.PullRequestID)
	if err != nil {
		return fmt.Errorf("failed to update pull request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("pull request not found")
	}

	return nil
}

// Exists проверяет существование Pull Request
func (r *prRepository) Exists(ctx context.Context, prID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check pull request existence: %w", err)
	}
	return exists, nil
}

// GetByReviewer возвращает PR, где пользователь назначен ревьювером
func (r *prRepository) GetByReviewer(ctx context.Context, reviewerID string) ([]*models.PullRequestShort, error) {
	query := `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		INNER JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
		WHERE prr.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests by reviewer: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var prs []*models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, fmt.Errorf("failed to scan pull request: %w", err)
		}
		prs = append(prs, &pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return prs, nil
}

// AssignReviewer назначает ревьювера на PR
func (r *prRepository) AssignReviewer(ctx context.Context, prID, reviewerID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := r.assignReviewerTx(ctx, tx, prID, reviewerID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// assignReviewerTx назначает ревьювера внутри транзакции
func (r *prRepository) assignReviewerTx(ctx context.Context, tx *sql.Tx, prID, reviewerID string) error {
	query := `
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
		ON CONFLICT (pull_request_id, reviewer_id) DO NOTHING
	`
	_, err := tx.ExecContext(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to assign reviewer: %w", err)
	}
	return nil
}

// RemoveReviewer удаляет ревьювера из PR
func (r *prRepository) RemoveReviewer(ctx context.Context, prID, reviewerID string) error {
	query := `
		DELETE FROM pr_reviewers
		WHERE pull_request_id = $1 AND reviewer_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, prID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to remove reviewer: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("reviewer assignment not found")
	}

	return nil
}

// GetReviewers возвращает список ревьюверов PR
func (r *prRepository) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	query := `
		SELECT reviewer_id
		FROM pr_reviewers
		WHERE pull_request_id = $1
		ORDER BY assigned_at
	`

	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var reviewers []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}
		reviewers = append(reviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return reviewers, nil
}

// IsReviewerAssigned проверяет, назначен ли ревьювер на PR
func (r *prRepository) IsReviewerAssigned(ctx context.Context, prID, reviewerID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM pr_reviewers
			WHERE pull_request_id = $1 AND reviewer_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, prID, reviewerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check reviewer assignment: %w", err)
	}

	return exists, nil
}
