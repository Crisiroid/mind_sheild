package service

import (
	"context"
	"errors"
	"fmt"

	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type MentalMustService struct {
	mentalMustRepo interfaces.MentalMustRepositoryInterface
}

func NewMentalMustService(mentalMustRepo interfaces.MentalMustRepositoryInterface) *MentalMustService {
	return &MentalMustService{
		mentalMustRepo: mentalMustRepo,
	}
}

func (s *MentalMustService) CreateMentalMust(ctx context.Context, req *schemas.MentalMustCreateRequest) (*schemas.MentalMustResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.MustText == "" {
		return nil, errors.New("متن باید الزامی است")
	}

	mentalMust := &models.MentalMust{
		UserID:    req.UserID,
		MustText:  req.MustText,
		DayNumber: req.DayNumber,
	}

	if err := s.mentalMustRepo.Create(ctx, mentalMust); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد باید ذهنی: %w", err)
	}

	return s.toMentalMustResponse(mentalMust), nil
}

func (s *MentalMustService) GetMentalMustById(ctx context.Context, id string) (*schemas.MentalMustResponse, error) {
	mentalMust, err := s.mentalMustRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("باید ذهنی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت باید ذهنی: %w", err)
	}

	return s.toMentalMustResponse(mentalMust), nil
}

func (s *MentalMustService) GetAllMentalMusts(ctx context.Context, req *schemas.MentalMustListRequest) (*schemas.MentalMustListResponse, error) {
	filterFunc := s.buildMentalMustFilters(req)

	mentalMusts, total, err := s.mentalMustRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت باید‌های ذهنی: %w", err)
	}

	responses := make([]schemas.MentalMustResponse, len(mentalMusts))
	for i, mentalMust := range mentalMusts {
		responses[i] = *s.toMentalMustResponse(&mentalMust)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.MentalMustListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.MentalMustResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *MentalMustService) UpdateMentalMust(ctx context.Context, id string, req *schemas.MentalMustUpdateRequest) (*schemas.MentalMustResponse, error) {
	mentalMust, err := s.mentalMustRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("باید ذهنی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت باید ذهنی: %w", err)
	}

	if req.IsReleased != nil {
		mentalMust.IsReleased = *req.IsReleased
	}

	if err := s.mentalMustRepo.Update(ctx, mentalMust); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی باید ذهنی: %w", err)
	}

	return s.toMentalMustResponse(mentalMust), nil
}

func (s *MentalMustService) DeleteMentalMust(ctx context.Context, id string) error {
	_, err := s.mentalMustRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("باید ذهنی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت باید ذهنی: %w", err)
	}

	if err := s.mentalMustRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف باید ذهنی: %w", err)
	}

	return nil
}

func (s *MentalMustService) buildMentalMustFilters(req *schemas.MentalMustListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.IsReleased != nil {
			db = db.Where("is_released = ?", *req.IsReleased)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *MentalMustService) toMentalMustResponse(mentalMust *models.MentalMust) *schemas.MentalMustResponse {
	return &schemas.MentalMustResponse{
		ID:           mentalMust.ID.String(),
		UserID:       mentalMust.UserID,
		MustText:     mentalMust.MustText,
		CreatedDate:  mentalMust.CreatedDate,
		IsReleased:   mentalMust.IsReleased,
		ReleasedDate: mentalMust.ReleasedDate,
		DayNumber:    mentalMust.DayNumber,
		CreatedAt:    mentalMust.CreatedAt,
		UpdatedAt:    mentalMust.UpdatedAt,
	}
}
