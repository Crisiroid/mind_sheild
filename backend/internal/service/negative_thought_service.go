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

type NegativeThoughtService struct {
	negativeThoughtRepo *repository.NegativeThoughtRepository
}

func NewNegativeThoughtService(negativeThoughtRepo *repository.NegativeThoughtRepository) *NegativeThoughtService {
	return &NegativeThoughtService{
		negativeThoughtRepo: negativeThoughtRepo,
	}
}

func (s *NegativeThoughtService) CreateNegativeThought(ctx context.Context, req *schemas.NegativeThoughtCreateRequest) (*schemas.NegativeThoughtResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.ThoughtText == "" {
		return nil, errors.New("متن فکر الزامی است")
	}

	thought := &models.NegativeThought{
		UserID:             req.UserID,
		ThoughtText:        req.ThoughtText,
		Situation:          req.Situation,
		CognitiveErrorType: req.CognitiveErrorType,
		ImpactLevel:        req.ImpactLevel,
		DayNumber:          req.DayNumber,
	}

	if err := s.negativeThoughtRepo.Create(ctx, thought); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد فکر منفی: %w", err)
	}

	return s.toNegativeThoughtResponse(thought), nil
}

func (s *NegativeThoughtService) GetNegativeThoughtById(ctx context.Context, id string) (*schemas.NegativeThoughtResponse, error) {
	thought, err := s.negativeThoughtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("فکر منفی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت فکر منفی: %w", err)
	}

	return s.toNegativeThoughtResponse(thought), nil
}

func (s *NegativeThoughtService) GetAllNegativeThoughts(ctx context.Context, req *schemas.NegativeThoughtListRequest) (*schemas.NegativeThoughtListResponse, error) {
	filterFunc := s.buildNegativeThoughtFilters(req)

	thoughts, total, err := s.negativeThoughtRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت افکار منفی: %w", err)
	}

	responses := make([]schemas.NegativeThoughtResponse, len(thoughts))
	for i, thought := range thoughts {
		responses[i] = *s.toNegativeThoughtResponse(&thought)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.NegativeThoughtListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.NegativeThoughtResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *NegativeThoughtService) UpdateNegativeThought(ctx context.Context, id string, req *schemas.NegativeThoughtUpdateRequest) (*schemas.NegativeThoughtResponse, error) {
	thought, err := s.negativeThoughtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("فکر منفی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت فکر منفی: %w", err)
	}

	if req.Situation != nil {
		thought.Situation = *req.Situation
	}

	if req.CognitiveErrorType != nil {
		thought.CognitiveErrorType = *req.CognitiveErrorType
	}

	if req.ImpactLevel != nil {
		thought.ImpactLevel = req.ImpactLevel
	}

	if err := s.negativeThoughtRepo.Update(ctx, thought); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی فکر منفی: %w", err)
	}

	return s.toNegativeThoughtResponse(thought), nil
}

func (s *NegativeThoughtService) DeleteNegativeThought(ctx context.Context, id string) error {
	_, err := s.negativeThoughtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("فکر منفی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت فکر منفی: %w", err)
	}

	if err := s.negativeThoughtRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف فکر منفی: %w", err)
	}

	return nil
}

func (s *NegativeThoughtService) buildNegativeThoughtFilters(req *schemas.NegativeThoughtListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.CognitiveErrorType != nil {
			db = db.Where("cognitive_error_type = ?", *req.CognitiveErrorType)
		}

		if req.ImpactMin != nil {
			db = db.Where("impact_level >= ?", *req.ImpactMin)
		}

		if req.ImpactMax != nil {
			db = db.Where("impact_level <= ?", *req.ImpactMax)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *NegativeThoughtService) toNegativeThoughtResponse(thought *models.NegativeThought) *schemas.NegativeThoughtResponse {
	return &schemas.NegativeThoughtResponse{
		ID:                 thought.ID.String(),
		UserID:             thought.UserID,
		ThoughtText:        thought.ThoughtText,
		Situation:          thought.Situation,
		CognitiveErrorType: thought.CognitiveErrorType,
		ImpactLevel:        thought.ImpactLevel,
		RecordedAt:         thought.RecordedAt,
		DayNumber:          thought.DayNumber,
		CreatedAt:          thought.CreatedAt,
		UpdatedAt:          thought.UpdatedAt,
	}
}
