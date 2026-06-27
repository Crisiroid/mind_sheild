package service

import (
	"context"
	"errors"
	"fmt"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type CognitiveGameService struct {
	cognitiveGameRepo *repository.CognitiveGameRepository
}

func NewCognitiveGameService(cognitiveGameRepo *repository.CognitiveGameRepository) *CognitiveGameService {
	return &CognitiveGameService{
		cognitiveGameRepo: cognitiveGameRepo,
	}
}

func (s *CognitiveGameService) CreateCognitiveGame(ctx context.Context, req *schemas.CognitiveGameCreateRequest) (*schemas.CognitiveGameResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	game := &models.CognitiveErrorGame{
		UserID:           req.UserID,
		ScenarioID:       req.ScenarioID,
		ScenarioType:     req.ScenarioType,
		Score:            req.Score,
		IsCorrect:        req.IsCorrect,
		TimeTakenSeconds: req.TimeTakenSeconds,
		DayNumber:        req.DayNumber,
	}

	if err := s.cognitiveGameRepo.Create(ctx, game); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد بازی شناختی: %w", err)
	}

	return s.toCognitiveGameResponse(game), nil
}

func (s *CognitiveGameService) GetCognitiveGameById(ctx context.Context, id string) (*schemas.CognitiveGameResponse, error) {
	game, err := s.cognitiveGameRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("بازی شناختی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت بازی شناختی: %w", err)
	}

	return s.toCognitiveGameResponse(game), nil
}

func (s *CognitiveGameService) GetAllCognitiveGames(ctx context.Context, req *schemas.CognitiveGameListRequest) (*schemas.CognitiveGameListResponse, error) {
	filterFunc := s.buildCognitiveGameFilters(req)

	games, total, err := s.cognitiveGameRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت بازی‌های شناختی: %w", err)
	}

	responses := make([]schemas.CognitiveGameResponse, len(games))
	for i, game := range games {
		responses[i] = *s.toCognitiveGameResponse(&game)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.CognitiveGameListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.CognitiveGameResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *CognitiveGameService) DeleteCognitiveGame(ctx context.Context, id string) error {
	_, err := s.cognitiveGameRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("بازی شناختی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت بازی شناختی: %w", err)
	}

	if err := s.cognitiveGameRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف بازی شناختی: %w", err)
	}

	return nil
}

func (s *CognitiveGameService) buildCognitiveGameFilters(req *schemas.CognitiveGameListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.ScenarioType != nil {
			db = db.Where("scenario_type = ?", *req.ScenarioType)
		}

		if req.IsCorrect != nil {
			db = db.Where("is_correct = ?", *req.IsCorrect)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *CognitiveGameService) toCognitiveGameResponse(game *models.CognitiveErrorGame) *schemas.CognitiveGameResponse {
	return &schemas.CognitiveGameResponse{
		ID:               game.ID.String(),
		UserID:           game.UserID,
		GameDate:         game.GameDate,
		ScenarioID:       game.ScenarioID,
		ScenarioType:     game.ScenarioType,
		Score:            game.Score,
		IsCorrect:        game.IsCorrect,
		TimeTakenSeconds: game.TimeTakenSeconds,
		DayNumber:        game.DayNumber,
		CreatedAt:        game.CreatedAt,
		UpdatedAt:        game.UpdatedAt,
	}
}
