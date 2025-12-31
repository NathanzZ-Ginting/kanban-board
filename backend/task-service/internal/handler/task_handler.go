package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/middleware"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/service"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.taskService.CreateTask(user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, dto.SuccessResponse{
		Message: "Task created successfully",
		Data:    task,
	})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskService.GetTaskByID(uint(id))
	if err != nil {
		h.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) GetTasksByBoard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid board ID")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	tasks, err := h.taskService.GetTasksByBoardID(uint(boardID), page, limit)
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTasksByColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	columnID, err := strconv.ParseUint(vars["columnId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid column ID")
		return
	}

	tasks, err := h.taskService.GetTasksByColumnID(uint(columnID))
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(uint(id), user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Task updated successfully",
		Data:    task,
	})
}

func (h *TaskHandler) MoveTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req dto.MoveTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.taskService.MoveTask(uint(id), &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Task moved successfully",
		Data:    task,
	})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := h.taskService.DeleteTask(uint(id), user.UserID); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Task deleted successfully",
	})
}

func (h *TaskHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	taskID, err := strconv.ParseUint(vars["taskId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	comment, err := h.taskService.AddComment(uint(taskID), user.UserID, &req)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusCreated, dto.SuccessResponse{
		Message: "Comment added successfully",
		Data:    comment,
	})
}

func (h *TaskHandler) GetTaskComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.ParseUint(vars["taskId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	comments, err := h.taskService.GetTaskComments(uint(taskID))
	if err != nil {
		h.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, comments)
}

func (h *TaskHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		h.errorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.ParseUint(vars["commentId"], 10, 32)
	if err != nil {
		h.errorResponse(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	if err := h.taskService.DeleteComment(uint(commentID), user.UserID); err != nil {
		h.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.jsonResponse(w, http.StatusOK, dto.SuccessResponse{
		Message: "Comment deleted successfully",
	})
}

func (h *TaskHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *TaskHandler) errorResponse(w http.ResponseWriter, status int, message string) {
	h.jsonResponse(w, status, dto.ErrorResponse{Error: message})
}
