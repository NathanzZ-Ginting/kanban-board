package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/config"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/handler"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/middleware"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/repository"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/service"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/pkg/database"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	log := logger.NewLogger(cfg.Environment)

	db, err := database.NewPostgresDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		// Check if error is "relation already exists" which is not critical
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatal("Failed to run database migrations", "error", err)
		}
		log.Warn("Database tables already exist, skipping migration", "error", err)
	}

	log.Info("Connected to Supabase PostgreSQL and migrations completed")

	userRepo := repository.NewUserRepository(db.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware(log))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"user-service"}`))
	}).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protected.HandleFunc("/profile", userHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", userHandler.UpdateProfile).Methods("PUT", "PATCH")
	protected.HandleFunc("/profile/password", userHandler.UpdatePassword).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

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
		log.Info("Starting User Service", "address", addr)
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
