package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/response"
	"github.com/zazaza5818/pr-reviewer-service/internal/service"
)

// PRHandler обрабатывает запросы к Pull Request
type PRHandler struct {
	service service.PullRequestService
}

// NewPRHandler создает новый обработчик Pull Request
func NewPRHandler(service service.PullRequestService) *PRHandler {
	return &PRHandler{service: service}
}

// CreatePR обрабатывает POST /pullRequest/create
func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "invalid request body")
		return
	}

	// Валидация
	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "pull_request_id, pull_request_name, and author_id are required")
		return
	}

	ctx := r.Context()
	pr, err := h.service.CreatePullRequest(ctx, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if errors.Is(err, service.ErrPRExists) {
			response.Error(w, http.StatusConflict, models.ErrPRExists, "PR id already exists")
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "author not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to create pull request")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

// MergePR обрабатывает POST /pullRequest/merge
func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "invalid request body")
		return
	}

	if req.PullRequestID == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "pull_request_id is required")
		return
	}

	ctx := r.Context()
	pr, err := h.service.MergePullRequest(ctx, req.PullRequestID)
	if err != nil {
		if errors.Is(err, service.ErrPRNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "pull request not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to merge pull request")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"pr": pr,
	})
}

// ReassignPR обрабатывает POST /pullRequest/reassign
func (h *PRHandler) ReassignPR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "invalid request body")
		return
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "pull_request_id and old_user_id are required")
		return
	}

	ctx := r.Context()
	pr, newReviewerID, err := h.service.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		if errors.Is(err, service.ErrPRNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "pull request not found")
			return
		}
		if errors.Is(err, service.ErrPRMerged) {
			response.Error(w, http.StatusConflict, models.ErrPRMerged, "cannot reassign on merged PR")
			return
		}
		if errors.Is(err, service.ErrReviewerNotFound) {
			response.Error(w, http.StatusConflict, models.ErrNotAssigned, "reviewer is not assigned to this PR")
			return
		}
		if errors.Is(err, service.ErrNoCandidate) {
			response.Error(w, http.StatusConflict, models.ErrNoCandidate, "no active replacement candidate in team")
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to reassign reviewer")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}
