package domain

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	OwnerID     uint           `gorm:"not null;index" json:"ownerId"`
	IsPublic    bool           `gorm:"default:false" json:"isPublic"`
	Color       string         `gorm:"size:20;default:#3b82f6" json:"color"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Columns     []Column       `gorm:"foreignKey:BoardID" json:"columns,omitempty"`
	Members     []BoardMember  `gorm:"foreignKey:BoardID" json:"members,omitempty"`
}

func (Board) TableName() string {
	return "boards"
}

type Column struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BoardID   uint           `gorm:"not null;index" json:"boardId"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Position  int            `gorm:"not null;default:0" json:"position"`
	Color     string         `gorm:"size:20;default:#e5e7eb" json:"color"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Column) TableName() string {
	return "columns"
}

type BoardMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BoardID   uint           `gorm:"not null;index" json:"boardId"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	Role      string         `gorm:"size:20;default:member" json:"role"` // owner, admin, member
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BoardMember) TableName() string {
	return "board_members"
}
