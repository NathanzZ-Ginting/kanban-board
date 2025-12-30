package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/kanban-monorepo/backend/auth-service/config"
	"github.com/yourusername/kanban-monorepo/backend/auth-service/internal/handler"
	"github.com/yourusername/kanban-monorepo/backend/auth-service/internal/middleware"
	"github.com/yourusername/kanban-monorepo/backend/auth-service/internal/repository"
	"github.com/yourusername/kanban-monorepo/backend/auth-service/internal/service"
	"github.com/yourusername/kanban-monorepo/backend/auth-service/pkg/database"
	"github.com/yourusername/kanban-monorepo/backend/auth-service/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.Environment)

	// Initialize MySQL database
	db, err := database.NewMySQLDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Run database migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to run database migrations", "error", err)
	}

	log.Info("Connected to MySQL and migrations completed")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.GetDB())

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Initialize router
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware(log))

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	authRoutes.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","database":"mysql"}`))
	}).Methods("GET")

	// Server configuration
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info("Starting Auth Service", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown server", "error", err)
	}

	log.Info("Server exited")
}
