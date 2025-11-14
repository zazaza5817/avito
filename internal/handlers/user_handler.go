package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/response"
	"github.com/zazaza5818/pr-reviewer-service/internal/service"
)

// UserHandler обрабатывает запросы к пользователям
type UserHandler struct {
	service service.UserService
}

// NewUserHandler создает новый обработчик пользователей
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// SetIsActive обрабатывает POST /users/setIsActive
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "invalid request body")
		return
	}

	if req.UserID == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "user_id is required")
		return
	}

	ctx := r.Context()
	user, err := h.service.SetUserActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to update user")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

// GetReviews обрабатывает GET /users/getReview?user_id=...
func (h *UserHandler) GetReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "user_id query parameter is required")
		return
	}

	ctx := r.Context()
	prs, err := h.service.GetUserReviews(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to get user reviews")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
