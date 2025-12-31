package service

import (
	"errors"

	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/domain"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/dto"
	"github.com/NathanzZ-Ginting/kanban-monorepo/backend/board-service/internal/repository"
)

type BoardService interface {
	CreateBoard(ownerID uint, req *dto.CreateBoardRequest) (*dto.BoardResponse, error)
	GetBoardByID(id uint) (*dto.BoardResponse, error)
	GetUserBoards(userID uint, page, limit int) (*dto.BoardsListResponse, error)
	UpdateBoard(id, userID uint, req *dto.UpdateBoardRequest) (*dto.BoardResponse, error)
	DeleteBoard(id, userID uint) error
	CreateColumn(boardID, userID uint, req *dto.CreateColumnRequest) (*dto.ColumnResponse, error)
	UpdateColumn(columnID, userID uint, req *dto.UpdateColumnRequest) (*dto.ColumnResponse, error)
	DeleteColumn(columnID, userID uint) error
	GetBoardColumns(boardID uint) ([]dto.ColumnResponse, error)
}

type boardService struct {
	boardRepo  repository.BoardRepository
	columnRepo repository.ColumnRepository
}

func NewBoardService(boardRepo repository.BoardRepository, columnRepo repository.ColumnRepository) BoardService {
	return &boardService{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
	}
}

func (s *boardService) CreateBoard(ownerID uint, req *dto.CreateBoardRequest) (*dto.BoardResponse, error) {
	if req.Name == "" {
		return nil, errors.New("board name is required")
	}

	color := req.Color
	if color == "" {
		color = "#3b82f6"
	}

	board := &domain.Board{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		IsPublic:    req.IsPublic,
		Color:       color,
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, err
	}

	// Create default columns
	defaultColumns := []domain.Column{
		{BoardID: board.ID, Name: "To Do", Position: 0, Color: "#e5e7eb"},
		{BoardID: board.ID, Name: "In Progress", Position: 1, Color: "#fbbf24"},
		{BoardID: board.ID, Name: "Done", Position: 2, Color: "#22c55e"},
	}

	for _, col := range defaultColumns {
		s.columnRepo.Create(&col)
	}

	return s.toBoardResponse(board), nil
}

func (s *boardService) GetBoardByID(id uint) (*dto.BoardResponse, error) {
	board, err := s.boardRepo.FindByIDWithColumns(id)
	if err != nil {
		return nil, errors.New("board not found")
	}

	return s.toBoardResponseWithColumns(board), nil
}

func (s *boardService) GetUserBoards(userID uint, page, limit int) (*dto.BoardsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	boards, total, err := s.boardRepo.FindUserBoards(userID, page, limit)
	if err != nil {
		return nil, err
	}

	var boardResponses []dto.BoardResponse
	for _, board := range boards {
		boardResponses = append(boardResponses, *s.toBoardResponse(&board))
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &dto.BoardsListResponse{
		Boards:     boardResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *boardService) UpdateBoard(id, userID uint, req *dto.UpdateBoardRequest) (*dto.BoardResponse, error) {
	board, err := s.boardRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("board not found")
	}

	if board.OwnerID != userID {
		return nil, errors.New("unauthorized to update this board")
	}

	if req.Name != "" {
		board.Name = req.Name
	}
	if req.Description != "" {
		board.Description = req.Description
	}
	if req.IsPublic != nil {
		board.IsPublic = *req.IsPublic
	}
	if req.Color != "" {
		board.Color = req.Color
	}

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}

	return s.toBoardResponse(board), nil
}

func (s *boardService) DeleteBoard(id, userID uint) error {
	board, err := s.boardRepo.FindByID(id)
	if err != nil {
		return errors.New("board not found")
	}

	if board.OwnerID != userID {
		return errors.New("unauthorized to delete this board")
	}

	return s.boardRepo.Delete(id)
}

func (s *boardService) CreateColumn(boardID, userID uint, req *dto.CreateColumnRequest) (*dto.ColumnResponse, error) {
	board, err := s.boardRepo.FindByID(boardID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	if board.OwnerID != userID {
		return nil, errors.New("unauthorized to add column to this board")
	}

	if req.Name == "" {
		return nil, errors.New("column name is required")
	}

	color := req.Color
	if color == "" {
		color = "#e5e7eb"
	}

	column := &domain.Column{
		BoardID:  boardID,
		Name:     req.Name,
		Position: req.Position,
		Color:    color,
	}

	if err := s.columnRepo.Create(column); err != nil {
		return nil, err
	}

	return s.toColumnResponse(column), nil
}

func (s *boardService) UpdateColumn(columnID, userID uint, req *dto.UpdateColumnRequest) (*dto.ColumnResponse, error) {
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return nil, errors.New("column not found")
	}

	board, err := s.boardRepo.FindByID(column.BoardID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	if board.OwnerID != userID {
		return nil, errors.New("unauthorized to update this column")
	}

	if req.Name != "" {
		column.Name = req.Name
	}
	if req.Position != nil {
		column.Position = *req.Position
	}
	if req.Color != "" {
		column.Color = req.Color
	}

	if err := s.columnRepo.Update(column); err != nil {
		return nil, err
	}

	return s.toColumnResponse(column), nil
}

func (s *boardService) DeleteColumn(columnID, userID uint) error {
	column, err := s.columnRepo.FindByID(columnID)
	if err != nil {
		return errors.New("column not found")
	}

	board, err := s.boardRepo.FindByID(column.BoardID)
	if err != nil {
		return errors.New("board not found")
	}

	if board.OwnerID != userID {
		return errors.New("unauthorized to delete this column")
	}

	return s.columnRepo.Delete(columnID)
}

func (s *boardService) GetBoardColumns(boardID uint) ([]dto.ColumnResponse, error) {
	columns, err := s.columnRepo.FindByBoardID(boardID)
	if err != nil {
		return nil, err
	}

	var columnResponses []dto.ColumnResponse
	for _, col := range columns {
		columnResponses = append(columnResponses, *s.toColumnResponse(&col))
	}

	return columnResponses, nil
}

func (s *boardService) toBoardResponse(board *domain.Board) *dto.BoardResponse {
	return &dto.BoardResponse{
		ID:          board.ID,
		Name:        board.Name,
		Description: board.Description,
		OwnerID:     board.OwnerID,
		IsPublic:    board.IsPublic,
		Color:       board.Color,
		CreatedAt:   board.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   board.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *boardService) toBoardResponseWithColumns(board *domain.Board) *dto.BoardResponse {
	resp := s.toBoardResponse(board)

	var columns []dto.ColumnResponse
	for _, col := range board.Columns {
		columns = append(columns, *s.toColumnResponse(&col))
	}
	resp.Columns = columns

	return resp
}

func (s *boardService) toColumnResponse(column *domain.Column) *dto.ColumnResponse {
	return &dto.ColumnResponse{
		ID:        column.ID,
		BoardID:   column.BoardID,
		Name:      column.Name,
		Position:  column.Position,
		Color:     column.Color,
		CreatedAt: column.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: column.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
