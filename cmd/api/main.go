// Package main является точкой входа в сервис назначения ревьюеров для Pull Request'ов.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/zazaza5818/pr-reviewer-service/internal/auth"
	"github.com/zazaza5818/pr-reviewer-service/internal/config"
	"github.com/zazaza5818/pr-reviewer-service/internal/database"
	"github.com/zazaza5818/pr-reviewer-service/internal/handlers"
	"github.com/zazaza5818/pr-reviewer-service/internal/middleware"
	"github.com/zazaza5818/pr-reviewer-service/internal/repository"
	"github.com/zazaza5818/pr-reviewer-service/internal/service"
)

func main() {
	// конфигурация сервиса
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	db, err := database.New(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	log.Println("Successfully connected to database")

	// Генерируем JWT токены для тестирования
	adminToken, err := auth.GenerateToken("admin-user-id", true)
	if err != nil {
		log.Fatalf("Failed to generate admin token: %v", err)
	}

	userToken, err := auth.GenerateToken("regular-user-id", false)
	if err != nil {
		log.Fatalf("Failed to generate user token: %v", err)
	}

	log.Println("=== JWT Tokens for testing ===")
	log.Printf("Admin JWT: %s", adminToken)
	log.Printf("User JWT: %s", userToken)
	log.Println("==============================")

	// Инициализируем репозитории
	teamRepo := repository.NewTeamRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	prRepo := repository.NewPullRequestRepository(db.DB)

	// Инициализируем сервисы
	teamService := service.NewTeamService(teamRepo, userRepo)
	userService := service.NewUserService(userRepo, prRepo)
	prService := service.NewPullRequestService(userRepo, prRepo)

	// Инициализируем обработчики
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	prHandler := handlers.NewPRHandler(prService)
	healthHandler := handlers.NewHealthHandler()

	// Настраиваем роутер
	router := mux.NewRouter()

	// Регистрируем маршруты
	router.HandleFunc("/health", healthHandler.Check).Methods("GET")

	// Team routes (требуют аутентификацию)
	router.Handle("/team/add", middleware.RequireAuth(http.HandlerFunc(teamHandler.CreateTeam))).Methods("POST")
	router.Handle("/team/get", middleware.RequireAuth(http.HandlerFunc(teamHandler.GetTeam))).Methods("GET")

	// User routes
	// setIsActive требует admin токен
	router.Handle("/users/setIsActive", middleware.RequireAdmin(http.HandlerFunc(userHandler.SetIsActive))).Methods("POST")
	// getReview требует обычную аутентификацию
	router.Handle("/users/getReview", middleware.RequireAuth(http.HandlerFunc(userHandler.GetReviews))).Methods("GET")

	// PR routes (требуют admin токен)
	router.Handle("/pullRequest/create", middleware.RequireAdmin(http.HandlerFunc(prHandler.CreatePR))).Methods("POST")
	router.Handle("/pullRequest/merge", middleware.RequireAdmin(http.HandlerFunc(prHandler.MergePR))).Methods("POST")
	router.Handle("/pullRequest/reassign", middleware.RequireAdmin(http.HandlerFunc(prHandler.ReassignPR))).Methods("POST")

	// Middleware для логирования
	router.Use(middleware.Logging)

	// Настраиваем HTTP сервер
	srv := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// сервер в отдельной горутине
	go func() {
		log.Printf("Starting server on %s", cfg.GetServerAddr())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
