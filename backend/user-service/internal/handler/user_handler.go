package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/middleware"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/user-service/internal/service"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userResponse, err := h.userService.GetUserByID(user.UserID)
	if err != nil {
		h.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, userResponse)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userResponse, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		h.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, userResponse)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, err := h.userService.GetAllUsers(page, limit)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, users)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userResponse, err := h.userService.UpdateProfile(user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Profile updated successfully",
		Data:    userResponse,
	})
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.UpdatePassword(user.UserID, &req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Password updated successfully",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		h.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) errorResponse(w http.ResponseWriter, status int, message string) {
	h.jsonResponse(w, status, dto.ErrorResponse{Error: message})
}
