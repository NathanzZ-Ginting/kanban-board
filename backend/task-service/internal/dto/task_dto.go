package dto

type CreateTaskRequest struct {
	BoardID     uint    `json:"boardId"`
	ColumnID    uint    `json:"columnId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	AssigneeID  *uint   `json:"assigneeId"`
	DueDate     *string `json:"dueDate"`
}

type UpdateTaskRequest struct {
	ColumnID    *uint   `json:"columnId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	Position    *int    `json:"position"`
	AssigneeID  *uint   `json:"assigneeId"`
	DueDate     *string `json:"dueDate"`
}

type MoveTaskRequest struct {
	ColumnID uint `json:"columnId"`
	Position int  `json:"position"`
}

type TaskResponse struct {
	ID          uint              `json:"id"`
	BoardID     uint              `json:"boardId"`
	ColumnID    uint              `json:"columnId"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Priority    string            `json:"priority"`
	Position    int               `json:"position"`
	AssigneeID  *uint             `json:"assigneeId"`
	CreatorID   uint              `json:"creatorId"`
	DueDate     *string           `json:"dueDate"`
	Labels      []LabelResponse   `json:"labels,omitempty"`
	Comments    []CommentResponse `json:"comments,omitempty"`
	CreatedAt   string            `json:"createdAt"`
	UpdatedAt   string            `json:"updatedAt"`
}

type TasksListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"totalPages"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	ID        uint   `json:"id"`
	TaskID    uint   `json:"taskId"`
	UserID    uint   `json:"userId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateLabelRequest struct {
	BoardID uint   `json:"boardId"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type LabelResponse struct {
	ID      uint   `json:"id"`
	BoardID uint   `json:"boardId"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
