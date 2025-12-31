package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/config"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/handler"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/middleware"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/repository"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/service"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/pkg/database"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/pkg/logger"

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

	boardRepo := repository.NewBoardRepository(db.GetDB())
	columnRepo := repository.NewColumnRepository(db.GetDB())
	boardService := service.NewBoardService(boardRepo, columnRepo)
	boardHandler := handler.NewBoardHandler(boardService)

	router := mux.NewRouter()

	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware(log))

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"board-service"}`))
	}).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Handle OPTIONS preflight requests (before auth middleware)
	api.HandleFunc("/boards", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/boards/{boardId}/columns", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")
	api.HandleFunc("/columns/{columnId}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Board routes
	protected.HandleFunc("/boards", boardHandler.CreateBoard).Methods("POST")
	protected.HandleFunc("/boards", boardHandler.GetUserBoards).Methods("GET")
	protected.HandleFunc("/boards/{id}", boardHandler.GetBoard).Methods("GET")
	protected.HandleFunc("/boards/{id}", boardHandler.UpdateBoard).Methods("PUT", "PATCH")
	protected.HandleFunc("/boards/{id}", boardHandler.DeleteBoard).Methods("DELETE")

	// Column routes
	protected.HandleFunc("/boards/{boardId}/columns", boardHandler.CreateColumn).Methods("POST")
	protected.HandleFunc("/boards/{boardId}/columns", boardHandler.GetBoardColumns).Methods("GET")
	protected.HandleFunc("/columns/{columnId}", boardHandler.UpdateColumn).Methods("PUT", "PATCH")
	protected.HandleFunc("/columns/{columnId}", boardHandler.DeleteColumn).Methods("DELETE")

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
		log.Info("Starting Board Service", "address", addr)
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
