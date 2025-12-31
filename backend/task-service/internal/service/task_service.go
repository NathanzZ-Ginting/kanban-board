package service

import (
	"errors"
	"time"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/domain"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/repository"
)

type TaskService interface {
	CreateTask(creatorID uint, req *dto.CreateTaskRequest) (*dto.TaskResponse, error)
	GetTaskByID(id uint) (*dto.TaskResponse, error)
	GetTasksByBoardID(boardID uint, page, limit int) (*dto.TasksListResponse, error)
	GetTasksByColumnID(columnID uint) ([]dto.TaskResponse, error)
	UpdateTask(id, userID uint, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error)
	MoveTask(id uint, req *dto.MoveTaskRequest) (*dto.TaskResponse, error)
	DeleteTask(id, userID uint) error
	AddComment(taskID, userID uint, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetTaskComments(taskID uint) ([]dto.CommentResponse, error)
	DeleteComment(commentID, userID uint) error
}

type taskService struct {
	taskRepo    repository.TaskRepository
	commentRepo repository.CommentRepository
}

func NewTaskService(taskRepo repository.TaskRepository, commentRepo repository.CommentRepository) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		commentRepo: commentRepo,
	}
}

func (s *taskService) CreateTask(creatorID uint, req *dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	if req.Title == "" {
		return nil, errors.New("task title is required")
	}

	if req.BoardID == 0 || req.ColumnID == 0 {
		return nil, errors.New("board ID and column ID are required")
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	task := &domain.Task{
		BoardID:     req.BoardID,
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		CreatedBy:   creatorID,
		AssigneeID:  req.AssigneeID,
	}

	if req.DueDate != nil {
		dueDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			task.DueDate = &dueDate
		}
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return s.toTaskResponse(task), nil
}

func (s *taskService) GetTaskByID(id uint) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByIDWithDetails(id)
	if err != nil {
		return nil, errors.New("task not found")
	}

	return s.toTaskResponseWithDetails(task), nil
}

func (s *taskService) GetTasksByBoardID(boardID uint, page, limit int) (*dto.TasksListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	tasks, total, err := s.taskRepo.FindByBoardID(boardID, page, limit)
	if err != nil {
		return nil, err
	}

	var taskResponses []dto.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, *s.toTaskResponse(&task))
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &dto.TasksListResponse{
		Tasks:      taskResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *taskService) GetTasksByColumnID(columnID uint) ([]dto.TaskResponse, error) {
	tasks, err := s.taskRepo.FindByColumnID(columnID)
	if err != nil {
		return nil, err
	}

	var taskResponses []dto.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, *s.toTaskResponse(&task))
	}

	return taskResponses, nil
}

func (s *taskService) UpdateTask(id, userID uint, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}

	if req.ColumnID != nil {
		task.ColumnID = *req.ColumnID
	}
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.Position != nil {
		task.Position = *req.Position
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		dueDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			task.DueDate = &dueDate
		}
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return s.toTaskResponse(task), nil
}

func (s *taskService) MoveTask(id uint, req *dto.MoveTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("task not found")
	}

	task.ColumnID = req.ColumnID
	task.Position = req.Position

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return s.toTaskResponse(task), nil
}

func (s *taskService) DeleteTask(id, userID uint) error {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return errors.New("task not found")
	}

	// Only creator can delete the task
	if task.CreatedBy != userID {
		return errors.New("unauthorized to delete this task")
	}

	return s.taskRepo.Delete(id)
}

func (s *taskService) AddComment(taskID, userID uint, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	if req.Content == "" {
		return nil, errors.New("comment content is required")
	}

	_, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	comment := &domain.Comment{
		TaskID:  taskID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment), nil
}

func (s *taskService) GetTaskComments(taskID uint) ([]dto.CommentResponse, error) {
	comments, err := s.commentRepo.FindByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	var commentResponses []dto.CommentResponse
	for _, comment := range comments {
		commentResponses = append(commentResponses, *s.toCommentResponse(&comment))
	}

	return commentResponses, nil
}

func (s *taskService) DeleteComment(commentID, userID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	if comment.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	return s.commentRepo.Delete(commentID)
}

func (s *taskService) toTaskResponse(task *domain.Task) *dto.TaskResponse {
	resp := &dto.TaskResponse{
		ID:          task.ID,
		BoardID:     task.BoardID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Position:    task.Position,
		AssigneeID:  task.AssigneeID,
		CreatorID:   task.CreatedBy,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if task.DueDate != nil {
		dueDate := task.DueDate.Format("2006-01-02")
		resp.DueDate = &dueDate
	}

	return resp
}

func (s *taskService) toTaskResponseWithDetails(task *domain.Task) *dto.TaskResponse {
	resp := s.toTaskResponse(task)

	var labels []dto.LabelResponse
	for _, label := range task.Labels {
		labels = append(labels, dto.LabelResponse{
			ID:      label.ID,
			BoardID: label.BoardID,
			Name:    label.Name,
			Color:   label.Color,
		})
	}
	resp.Labels = labels

	var comments []dto.CommentResponse
	for _, comment := range task.Comments {
		comments = append(comments, *s.toCommentResponse(&comment))
	}
	resp.Comments = comments

	return resp
}

func (s *taskService) toCommentResponse(comment *domain.Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:        comment.ID,
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
