package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/middleware"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/service"
	"github.com/gorilla/mux"
)

type BoardHandler struct {
	boardService service.BoardService
}

func NewBoardHandler(boardService service.BoardService) *BoardHandler {
	return &BoardHandler{
		boardService: boardService,
	}
}

func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	board, err := h.boardService.CreateBoard(user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, dto.SuccessResponse{
		Message: "Board created successfully",
		Data:    board,
	})
}

func (h *BoardHandler) GetBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	board, err := h.boardService.GetBoardByID(uint(id))
	if err != nil {
		h.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, board)
}

func (h *BoardHandler) GetUserBoards(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	boards, err := h.boardService.GetUserBoards(user.UserID, page, limit)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, boards)
}

func (h *BoardHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var req dto.UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	board, err := h.boardService.UpdateBoard(uint(id), user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Board updated successfully",
		Data:    board,
	})
}

func (h *BoardHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	if err := h.boardService.DeleteBoard(uint(id), user.UserID); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Board deleted successfully",
	})
}

func (h *BoardHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	var req dto.CreateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	column, err := h.boardService.CreateColumn(uint(boardID), user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, dto.SuccessResponse{
		Message: "Column created successfully",
		Data:    column,
	})
}

func (h *BoardHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	columnID, err := strconv.ParseUint(vars["columnId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid column ID")
		return
	}

	var req dto.UpdateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	column, err := h.boardService.UpdateColumn(uint(columnID), user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Column updated successfully",
		Data:    column,
	})
}

func (h *BoardHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	columnID, err := strconv.ParseUint(vars["columnId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid column ID")
		return
	}

	if err := h.boardService.DeleteColumn(uint(columnID), user.UserID); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Column deleted successfully",
	})
}

func (h *BoardHandler) GetBoardColumns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	columns, err := h.boardService.GetBoardColumns(uint(boardID))
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, columns)
}

func (h *BoardHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *BoardHandler) errorResponse(w http.ResponseWriter, status int, message string) {
	h.jsonResponse(w, status, dto.ErrorResponse{Error: message})
}
