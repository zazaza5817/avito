// Package response предоставляет утилиты для отправки HTTP ответов.
package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/zazaza5818/pr-reviewer-service/internal/models"
)

// JSON отправляет JSON ответ
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Error отправляет ответ с ошибкой
func Error(w http.ResponseWriter, statusCode int, code models.ErrorCode, message string) {
	errorResponse := models.NewErrorResponse(code, message)
	JSON(w, statusCode, errorResponse)
}
