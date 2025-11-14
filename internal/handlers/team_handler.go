package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/response"
	"github.com/zazaza5818/pr-reviewer-service/internal/service"
)

// TeamHandler обрабатывает запросы к командам
type TeamHandler struct {
	service service.TeamService
}

// NewTeamHandler создает новый обработчик команд
func NewTeamHandler(service service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

// CreateTeam обрабатывает POST /team/add
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "invalid request body")
		return
	}

	// Валидация
	if team.TeamName == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "team_name is required")
		return
	}

	if len(team.Members) == 0 {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "members are required")
		return
	}

	// Валидация членов команды
	for _, member := range team.Members {
		if member.UserID == "" || member.Username == "" {
			response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "user_id and username are required for all members")
			return
		}
	}

	ctx := r.Context()
	if err := h.service.CreateTeam(ctx, &team); err != nil {
		if errors.Is(err, service.ErrTeamExists) {
			response.Error(w, http.StatusBadRequest, models.ErrTeamExists, "team_name already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to create team")
		return
	}

	// Получаем созданную команду для ответа
	createdTeam, err := h.service.GetTeam(ctx, team.TeamName)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to retrieve created team")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"team": createdTeam,
	})
}

// GetTeam обрабатывает GET /team/get?team_name=...
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		response.Error(w, http.StatusBadRequest, models.ErrBadRequest, "team_name query parameter is required")
		return
	}

	ctx := r.Context()
	team, err := h.service.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {
			response.Error(w, http.StatusNotFound, models.ErrNotFound, "team not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, models.ErrInternal, "failed to get team")
		return
	}

	response.JSON(w, http.StatusOK, team)
}
