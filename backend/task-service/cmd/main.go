package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/config"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/handler"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/middleware"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/repository"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/service"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/pkg/database"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger(cfg.Environment)

	db, err := database.NewMySQLDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to run database migrations", "error", err)
	}

	log.Info("Connected to MySQL and migrations completed")

	taskRepo := repository.NewTaskRepository(db.GetDB())
	commentRepo := repository.NewCommentRepository(db.GetDB())
	taskService := service.NewTaskService(taskRepo, commentRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware(log))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"task-service"}`))
	}).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Handle OPTIONS preflight requests (before auth middleware)
	api.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/tasks/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/boards/{boardId}/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/columns/{columnId}/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/tasks/{taskId}/comments", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/comments/{commentId}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Task routes
	protected.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT", "PATCH")
	protected.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
	protected.HandleFunc("/tasks/{id}/move", taskHandler.MoveTask).Methods("PUT", "PATCH")

	// Board tasks
	protected.HandleFunc("/boards/{boardId}/tasks", taskHandler.GetTasksByBoard).Methods("GET")

	// Column tasks
	protected.HandleFunc("/columns/{columnId}/tasks", taskHandler.GetTasksByColumn).Methods("GET")

	// Comment routes
	protected.HandleFunc("/tasks/{taskId}/comments", taskHandler.AddComment).Methods("POST")
	protected.HandleFunc("/tasks/{taskId}/comments", taskHandler.GetTaskComments).Methods("GET")
	protected.HandleFunc("/comments/{commentId}", taskHandler.DeleteComment).Methods("DELETE")

	// Server setup
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Starting Task Service", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited properly")
}
