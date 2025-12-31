package repository

import (
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/task-service/internal/domain"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id uint) (*domain.Task, error)
	FindByIDWithDetails(id uint) (*domain.Task, error)
	FindByColumnID(columnID uint) ([]domain.Task, error)
	FindByBoardID(boardID uint, page, limit int) ([]domain.Task, int64, error)
	Update(task *domain.Task) error
	Delete(id uint) error
	UpdatePositions(tasks []domain.Task) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uint) (*domain.Task, error) {
	var task domain.Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByIDWithDetails(id uint) (*domain.Task, error) {
	var task domain.Task
	if err := r.db.Preload("Labels").Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindByColumnID(columnID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	if err := r.db.Where("column_id = ?", columnID).Order("position ASC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) FindByBoardID(boardID uint, page, limit int) ([]domain.Task, int64, error) {
	var tasks []domain.Task
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Task{}).Where("board_id = ?", boardID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Labels").Where("board_id = ?", boardID).
		Offset(offset).Limit(limit).Order("position ASC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Task{}, id).Error
}

func (r *taskRepository) UpdatePositions(tasks []domain.Task) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, task := range tasks {
			if err := tx.Model(&domain.Task{}).Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"position":  task.Position,
					"column_id": task.ColumnID,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
