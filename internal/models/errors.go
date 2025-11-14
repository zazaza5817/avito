// Package models определяет структуры данных и типы ошибок, используемые в приложении.
package models

// ErrorCode представляет код ошибки
type ErrorCode string

// Коды ошибок API
const (
	ErrTeamExists   ErrorCode = "TEAM_EXISTS"
	ErrPRExists     ErrorCode = "PR_EXISTS"
	ErrPRMerged     ErrorCode = "PR_MERGED"
	ErrNotAssigned  ErrorCode = "NOT_ASSIGNED"
	ErrNoCandidate  ErrorCode = "NO_CANDIDATE"
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
)

// ErrorDetail представляет детали ошибки
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// NewErrorResponse создает новый ErrorResponse
func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
