package domain

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	BoardID     uint           `gorm:"not null;index" json:"boardId"`
	ColumnID    uint           `gorm:"not null;index" json:"columnId"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Priority    string         `gorm:"size:20;default:medium" json:"priority"` // low, medium, high, urgent
	Position    int            `gorm:"not null;default:0" json:"position"`
	AssigneeID  *uint          `gorm:"index" json:"assigneeId"`
	CreatedBy   uint          `gorm:"column:created_by" json:"createdBy"`
	DueDate     *time.Time     `json:"dueDate"`
	Labels      []Label        `gorm:"many2many:task_labels;" json:"labels,omitempty"`
	Comments    []Comment      `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Task) TableName() string {
	return "tasks"
}

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TaskID    uint           `gorm:"not null;index" json:"taskId"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Comment) TableName() string {
	return "comments"
}

type Label struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BoardID   uint           `gorm:"not null;index" json:"boardId"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Color     string         `gorm:"size:20;default:#3b82f6" json:"color"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Label) TableName() string {
	return "labels"
}
