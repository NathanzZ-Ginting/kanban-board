package dto

type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
	Color       string `json:"color"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"isPublic"`
	Color       string `json:"color"`
}

type BoardResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	OwnerID     uint             `json:"ownerId"`
	IsPublic    bool             `json:"isPublic"`
	Color       string           `json:"color"`
	Columns     []ColumnResponse `json:"columns,omitempty"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

type BoardsListResponse struct {
	Boards     []BoardResponse `json:"boards"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"totalPages"`
}

type CreateColumnRequest struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
	Color    string `json:"color"`
}

type UpdateColumnRequest struct {
	Name     string `json:"name"`
	Position *int   `json:"position"`
	Color    string `json:"color"`
}

type ColumnResponse struct {
	ID        uint   `json:"id"`
	BoardID   uint   `json:"boardId"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type AddMemberRequest struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
}

type MemberResponse struct {
	ID        uint   `json:"id"`
	BoardID   uint   `json:"boardId"`
	UserID    uint   `json:"userId"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
