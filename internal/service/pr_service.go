package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/repository"
)

// PullRequestService определяет интерфейс для работы с Pull Request
type PullRequestService interface {
	CreatePullRequest(ctx context.Context, prID, prName, authorID string) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*models.PullRequest, string, error)
}

// pullRequestService реализует PullRequestService
type pullRequestService struct {
	userRepo repository.UserRepository
	prRepo   repository.PullRequestRepository
	rand     *rand.Rand
}

// NewPullRequestService создает новый сервис для работы с Pull Request
func NewPullRequestService(
	userRepo repository.UserRepository,
	prRepo repository.PullRequestRepository,
) PullRequestService {
	return &pullRequestService{
		userRepo: userRepo,
		prRepo:   prRepo,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreatePullRequest создает PR и автоматически назначает ревьюверов
func (s *pullRequestService) CreatePullRequest(ctx context.Context, prID, prName, authorID string) (*models.PullRequest, error) {
	// Проверяем существование PR
	exists, err := s.prRepo.Exists(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to check PR existence: %w", err)
	}
	if exists {
		return nil, ErrPRExists
	}

	// Получаем автора
	author, err := s.userRepo.Get(ctx, authorID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Получаем активных участников команды автора (исключая самого автора)
	candidates, err := s.userRepo.GetActiveTeammates(ctx, author.TeamName, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team candidates: %w", err)
	}

	// Выбираем до 2 случайных ревьюверов
	reviewers := s.selectRandomReviewers(candidates, 2)

	// Создаем PR
	pr := &models.PullRequest{
		PullRequestID:     prID,
		PullRequestName:   prName,
		AuthorID:          authorID,
		Status:            models.StatusOpen,
		AssignedReviewers: reviewers,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	// Получаем созданный PR с полными данными
	createdPR, err := s.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created PR: %w", err)
	}

	return createdPR, nil
}

// MergePullRequest помечает PR как MERGED (идемпотентная операция)
func (s *pullRequestService) MergePullRequest(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := s.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, ErrPRNotFound
	}

	// Если уже merged, просто возвращаем PR (идемпотентность)
	if pr.Status == models.StatusMerged {
		return pr, nil
	}

	// Обновляем статус
	pr.Status = models.StatusMerged
	now := time.Now()
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, pr); err != nil {
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	return pr, nil
}

// ReassignReviewer переназначает ревьювера
func (s *pullRequestService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*models.PullRequest, string, error) {
	// Получаем PR
	pr, err := s.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, "", ErrPRNotFound
	}

	// Проверяем, что PR не merged
	if pr.Status == models.StatusMerged {
		return nil, "", ErrPRMerged
	}

	// Проверяем, что oldReviewerID назначен на этот PR
	isAssigned, err := s.prRepo.IsReviewerAssigned(ctx, prID, oldReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check reviewer assignment: %w", err)
	}
	if !isAssigned {
		return nil, "", ErrReviewerNotFound
	}

	// Получаем старого ревьювера для определения его команды
	oldReviewer, err := s.userRepo.Get(ctx, oldReviewerID)
	if err != nil {
		return nil, "", ErrUserNotFound
	}

	// Получаем активных участников команды старого ревьювера
	// Исключаем автора PR и текущих ревьюверов
	candidates, err := s.userRepo.GetActiveTeammates(ctx, oldReviewer.TeamName, "")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get team candidates: %w", err)
	}

	// Фильтруем кандидатов: исключаем автора и уже назначенных ревьюверов
	var filteredCandidates []*models.User
	for _, candidate := range candidates {
		if candidate.UserID == pr.AuthorID {
			continue
		}
		alreadyAssigned := false
		for _, reviewerID := range pr.AssignedReviewers {
			if candidate.UserID == reviewerID {
				alreadyAssigned = true
				break
			}
		}
		if !alreadyAssigned {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	// Проверяем наличие кандидатов
	if len(filteredCandidates) == 0 {
		return nil, "", ErrNoCandidate
	}

	// Выбираем случайного кандидата
	newReviewer := filteredCandidates[s.rand.Intn(len(filteredCandidates))]

	// Удаляем старого ревьювера
	if err := s.prRepo.RemoveReviewer(ctx, prID, oldReviewerID); err != nil {
		return nil, "", fmt.Errorf("failed to remove old reviewer: %w", err)
	}

	// Назначаем нового ревьювера
	if err := s.prRepo.AssignReviewer(ctx, prID, newReviewer.UserID); err != nil {
		return nil, "", fmt.Errorf("failed to assign new reviewer: %w", err)
	}

	// Получаем обновленный PR
	updatedPR, err := s.prRepo.Get(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get updated PR: %w", err)
	}

	return updatedPR, newReviewer.UserID, nil
}

// selectRandomReviewers выбирает случайных ревьюверов из списка кандидатов
func (s *pullRequestService) selectRandomReviewers(candidates []*models.User, maxCount int) []string {
	count := len(candidates)
	if count > maxCount {
		count = maxCount
	}

	if count == 0 {
		return []string{}
	}

	// Создаем копию индексов
	indices := make([]int, len(candidates))
	for i := range indices {
		indices[i] = i
	}

	// Перемешиваем индексы
	s.rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// Выбираем первые count элементов
	reviewers := make([]string, count)
	for i := 0; i < count; i++ {
		reviewers[i] = candidates[indices[i]].UserID
	}

	return reviewers
}
