package repository

import (
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/domain"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *domain.Comment) error
	FindByID(id uint) (*domain.Comment, error)
	FindByTaskID(taskID uint) ([]domain.Comment, error)
	Update(comment *domain.Comment) error
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByTaskID(taskID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	if err := r.db.Where("task_id = ?", taskID).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Update(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}
