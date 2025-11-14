// Package middleware предоставляет HTTP middleware для аутентификации и логирования.
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/zazaza5818/pr-reviewer-service/internal/auth"
	"github.com/zazaza5818/pr-reviewer-service/internal/models"
	"github.com/zazaza5818/pr-reviewer-service/internal/response"
)

type contextKey string

const (
	// UserIDKey - ключ контекста для хранения ID пользователя
	UserIDKey contextKey = "user_id"
	// IsAdminKey - ключ контекста для хранения статуса администратора
	IsAdminKey contextKey = "is_admin"
)

// RequireAuth проверяет JWT токен и сохраняет claims в контексте
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, models.ErrUnauthorized, "missing authorization header")
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Error(w, http.StatusUnauthorized, models.ErrUnauthorized, "invalid authorization format, expected 'Bearer <token>'")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		if tokenString == "" {
			response.Error(w, http.StatusUnauthorized, models.ErrUnauthorized, "missing bearer token")
			return
		}

		// Валидируем JWT токен
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, models.ErrUnauthorized, "invalid or expired token")
			return
		}

		// Сохраняем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin проверяет что пользователь является администратором
func RequireAdmin(next http.Handler) http.Handler {
	return RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
		if !ok || !isAdmin {
			response.Error(w, http.StatusForbidden, models.ErrUnauthorized, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	}))
}

// Logging предоставляет middleware для логирования HTTP запросов.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.RequestURI,
			rw.statusCode,
			time.Since(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
