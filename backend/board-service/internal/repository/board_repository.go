package repository

import (
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/domain"
	"gorm.io/gorm"
)

type BoardRepository interface {
	Create(board *domain.Board) error
	FindByID(id uint) (*domain.Board, error)
	FindByIDWithColumns(id uint) (*domain.Board, error)
	FindByOwnerID(ownerID uint, page, limit int) ([]domain.Board, int64, error)
	FindUserBoards(userID uint, page, limit int) ([]domain.Board, int64, error)
	Update(board *domain.Board) error
	Delete(id uint) error
}

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(board *domain.Board) error {
	return r.db.Create(board).Error
}

func (r *boardRepository) FindByID(id uint) (*domain.Board, error) {
	var board domain.Board
	if err := r.db.First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) FindByIDWithColumns(id uint) (*domain.Board, error) {
	var board domain.Board
	if err := r.db.Preload("Columns", func(db *gorm.DB) *gorm.DB {
		return db.Order("position ASC")
	}).First(&board, id).Error; err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) FindByOwnerID(ownerID uint, page, limit int) ([]domain.Board, int64, error) {
	var boards []domain.Board
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&domain.Board{}).Where("owner_id = ?", ownerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("owner_id = ?", ownerID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (r *boardRepository) FindUserBoards(userID uint, page, limit int) ([]domain.Board, int64, error) {
	var boards []domain.Board
	var total int64

	offset := (page - 1) * limit

	// Find boards where user is owner or member
	subQuery := r.db.Model(&domain.BoardMember{}).Select("board_id").Where("user_id = ?", userID)

	countQuery := r.db.Model(&domain.Board{}).Where("owner_id = ? OR id IN (?)", userID, subQuery)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("owner_id = ? OR id IN (?)", userID, subQuery).
		Offset(offset).Limit(limit).Order("created_at DESC").Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

func (r *boardRepository) Update(board *domain.Board) error {
	return r.db.Save(board).Error
}

func (r *boardRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Board{}, id).Error
}
