package repository

import (
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/domain"
	"gorm.io/gorm"
)

type ColumnRepository interface {
	Create(column *domain.Column) error
	FindByID(id uint) (*domain.Column, error)
	FindByBoardID(boardID uint) ([]domain.Column, error)
	Update(column *domain.Column) error
	Delete(id uint) error
	UpdatePositions(columns []domain.Column) error
}

type columnRepository struct {
	db *gorm.DB
}

func NewColumnRepository(db *gorm.DB) ColumnRepository {
	return &columnRepository{db: db}
}

func (r *columnRepository) Create(column *domain.Column) error {
	return r.db.Create(column).Error
}

func (r *columnRepository) FindByID(id uint) (*domain.Column, error) {
	var column domain.Column
	if err := r.db.First(&column, id).Error; err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *columnRepository) FindByBoardID(boardID uint) ([]domain.Column, error) {
	var columns []domain.Column
	if err := r.db.Where("board_id = ?", boardID).Order("position ASC").Find(&columns).Error; err != nil {
		return nil, err
	}
	return columns, nil
}

func (r *columnRepository) Update(column *domain.Column) error {
	return r.db.Save(column).Error
}

func (r *columnRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Column{}, id).Error
}

func (r *columnRepository) UpdatePositions(columns []domain.Column) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, col := range columns {
			if err := tx.Model(&domain.Column{}).Where("id = ?", col.ID).Update("position", col.Position).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
