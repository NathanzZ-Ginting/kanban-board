package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate request
	if validationErrors := h.authService.ValidateRegisterRequest(&req); len(validationErrors) > 0 {
		sendError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	// Register user
	response, err := h.authService.Register(&req)
	if err != nil {
		switch err {
		case service.ErrEmailExists:
			sendError(w, http.StatusConflict, "Email already exists", nil)
		case service.ErrUsernameExists:
			sendError(w, http.StatusConflict, "Username already exists", nil)
		default:
			sendError(w, http.StatusInternalServerError, "Failed to register user", nil)
		}
		return
	}

	sendSuccess(w, http.StatusCreated, "User registered successfully", response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate request
	if validationErrors := h.authService.ValidateLoginRequest(&req); len(validationErrors) > 0 {
		sendError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	// Login user
	response, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			sendError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		default:
			sendError(w, http.StatusInternalServerError, "Failed to login", nil)
		}
		return
	}

	sendSuccess(w, http.StatusOK, "Login successful", response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		sendError(w, http.StatusBadRequest, "Refresh token is required", nil)
		return
	}

	response, err := h.authService.RefreshToken(&req)
	if err != nil {
		switch err {
		case service.ErrInvalidToken:
			sendError(w, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
		case service.ErrUserNotFound:
			sendError(w, http.StatusNotFound, "User not found", nil)
		default:
			sendError(w, http.StatusInternalServerError, "Failed to refresh token", nil)
		}
		return
	}

	sendSuccess(w, http.StatusOK, "Token refreshed successfully", response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if req.RefreshToken == "" {
		sendError(w, http.StatusBadRequest, "Refresh token is required", nil)
		return
	}

	if err := h.authService.Logout(&req); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to logout", nil)
		return
	}

	sendSuccess(w, http.StatusOK, "Logged out successfully", nil)
}

// Helper functions
func sendError(w http.ResponseWriter, status int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error:   message,
		Details: details,
	})
}

func sendSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	} else {
		json.NewEncoder(w).Encode(dto.SuccessResponse{
			Message: message,
		})
	}
}
