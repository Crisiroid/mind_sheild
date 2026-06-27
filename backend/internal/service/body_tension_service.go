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

type BodyTensionService struct {
	bodyTensionRepo *repository.BodyTensionRepository
}

func NewBodyTensionService(bodyTensionRepo *repository.BodyTensionRepository) *BodyTensionService {
	return &BodyTensionService{
		bodyTensionRepo: bodyTensionRepo,
	}
}

func (s *BodyTensionService) CreateBodyTension(ctx context.Context, req *schemas.BodyTensionCreateRequest) (*schemas.BodyTensionResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.BodyRegions == "" {
		return nil, errors.New("نواحی بدنی الزامی است")
	}

	bodyTension := &models.BodyTensionMap{
		UserID:           req.UserID,
		BodyRegions:      req.BodyRegions,
		OverallIntensity: req.OverallIntensity,
		SeverityColor:    req.SeverityColor,
		Notes:            req.Notes,
		DayNumber:        req.DayNumber,
	}

	if err := s.bodyTensionRepo.Create(ctx, bodyTension); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد نقشه تنش بدنی: %w", err)
	}

	return s.toBodyTensionResponse(bodyTension), nil
}

func (s *BodyTensionService) GetBodyTensionById(ctx context.Context, id string) (*schemas.BodyTensionResponse, error) {
	bodyTension, err := s.bodyTensionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نقشه تنش بدنی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت نقشه تنش بدنی: %w", err)
	}

	return s.toBodyTensionResponse(bodyTension), nil
}

func (s *BodyTensionService) GetAllBodyTensions(ctx context.Context, req *schemas.BodyTensionListRequest) (*schemas.BodyTensionListResponse, error) {
	filterFunc := s.buildBodyTensionFilters(req)

	bodyTensions, total, err := s.bodyTensionRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت نقشه‌های تنش بدنی: %w", err)
	}

	responses := make([]schemas.BodyTensionResponse, len(bodyTensions))
	for i, bodyTension := range bodyTensions {
		responses[i] = *s.toBodyTensionResponse(&bodyTension)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.BodyTensionListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.BodyTensionResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *BodyTensionService) UpdateBodyTension(ctx context.Context, id string, req *schemas.BodyTensionUpdateRequest) (*schemas.BodyTensionResponse, error) {
	bodyTension, err := s.bodyTensionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نقشه تنش بدنی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت نقشه تنش بدنی: %w", err)
	}

	if req.BodyRegions != nil {
		bodyTension.BodyRegions = *req.BodyRegions
	}

	if req.OverallIntensity != nil {
		bodyTension.OverallIntensity = req.OverallIntensity
	}

	if req.Notes != nil {
		bodyTension.Notes = *req.Notes
	}

	if err := s.bodyTensionRepo.Update(ctx, bodyTension); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی نقشه تنش بدنی: %w", err)
	}

	return s.toBodyTensionResponse(bodyTension), nil
}

func (s *BodyTensionService) DeleteBodyTension(ctx context.Context, id string) error {
	_, err := s.bodyTensionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("نقشه تنش بدنی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت نقشه تنش بدنی: %w", err)
	}

	if err := s.bodyTensionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف نقشه تنش بدنی: %w", err)
	}

	return nil
}

func (s *BodyTensionService) buildBodyTensionFilters(req *schemas.BodyTensionListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.IntensityMin != nil {
			db = db.Where("overall_intensity >= ?", *req.IntensityMin)
		}

		if req.IntensityMax != nil {
			db = db.Where("overall_intensity <= ?", *req.IntensityMax)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *BodyTensionService) toBodyTensionResponse(bodyTension *models.BodyTensionMap) *schemas.BodyTensionResponse {
	return &schemas.BodyTensionResponse{
		ID:               bodyTension.ID.String(),
		UserID:           bodyTension.UserID,
		MappingDate:      bodyTension.MappingDate,
		BodyRegions:      bodyTension.BodyRegions,
		OverallIntensity: bodyTension.OverallIntensity,
		SeverityColor:    bodyTension.SeverityColor,
		Notes:            bodyTension.Notes,
		DayNumber:        bodyTension.DayNumber,
		CreatedAt:        bodyTension.CreatedAt,
		UpdatedAt:        bodyTension.UpdatedAt,
	}
}
