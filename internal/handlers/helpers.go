// Package handlers содержит HTTP обработчики запросов для API эндпоинтов.
package handlers

import (
	"net/http"

	"github.com/zazaza5818/pr-reviewer-service/internal/response"
)

// HealthHandler обрабатывает проверку здоровья сервиса
type HealthHandler struct{}

// NewHealthHandler создает новый обработчик health check
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check обрабатывает GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
