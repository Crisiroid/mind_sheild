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

type SkyThoughtService struct {
	skyThoughtRepo interfaces.SkyThoughtRepositoryInterface
}

func NewSkyThoughtService(skyThoughtRepo interfaces.SkyThoughtRepositoryInterface) *SkyThoughtService {
	return &SkyThoughtService{
		skyThoughtRepo: skyThoughtRepo,
	}
}

func (s *SkyThoughtService) CreateSkyThought(ctx context.Context, req *schemas.SkyThoughtCreateRequest) (*schemas.SkyThoughtResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.ThoughtText == "" {
		return nil, errors.New("متن فکر الزامی است")
	}

	skyThought := &models.SkyThought{
		UserID:      req.UserID,
		ThoughtText: req.ThoughtText,
		DayNumber:   req.DayNumber,
	}

	if err := s.skyThoughtRepo.Create(ctx, skyThought); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد فکر آسمانی: %w", err)
	}

	return s.toSkyThoughtResponse(skyThought), nil
}

func (s *SkyThoughtService) GetSkyThoughtById(ctx context.Context, id string) (*schemas.SkyThoughtResponse, error) {
	skyThought, err := s.skyThoughtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("فکر آسمانی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت فکر آسمانی: %w", err)
	}

	return s.toSkyThoughtResponse(skyThought), nil
}

func (s *SkyThoughtService) GetAllSkyThoughts(ctx context.Context, req *schemas.SkyThoughtListRequest) (*schemas.SkyThoughtListResponse, error) {
	filterFunc := s.buildSkyThoughtFilters(req)

	skyThoughts, total, err := s.skyThoughtRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت افکار آسمانی: %w", err)
	}

	responses := make([]schemas.SkyThoughtResponse, len(skyThoughts))
	for i, skyThought := range skyThoughts {
		responses[i] = *s.toSkyThoughtResponse(&skyThought)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.SkyThoughtListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.SkyThoughtResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *SkyThoughtService) UpdateSkyThought(ctx context.Context, id string, req *schemas.SkyThoughtUpdateRequest) (*schemas.SkyThoughtResponse, error) {
	skyThought, err := s.skyThoughtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("فکر آسمانی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت فکر آسمانی: %w", err)
	}

	if req.CloudSwiped != nil {
		skyThought.CloudSwiped = *req.CloudSwiped
	}

	if err := s.skyThoughtRepo.Update(ctx, skyThought); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی فکر آسمانی: %w", err)
	}

	return s.toSkyThoughtResponse(skyThought), nil
}

func (s *SkyThoughtService) DeleteSkyThought(ctx context.Context, id string) error {
	_, err := s.skyThoughtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("فکر آسمانی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت فکر آسمانی: %w", err)
	}

	if err := s.skyThoughtRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف فکر آسمانی: %w", err)
	}

	return nil
}

func (s *SkyThoughtService) buildSkyThoughtFilters(req *schemas.SkyThoughtListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.CloudSwiped != nil {
			db = db.Where("cloud_swiped = ?", *req.CloudSwiped)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *SkyThoughtService) toSkyThoughtResponse(skyThought *models.SkyThought) *schemas.SkyThoughtResponse {
	return &schemas.SkyThoughtResponse{
		ID:          skyThought.ID.String(),
		UserID:      skyThought.UserID,
		ThoughtText: skyThought.ThoughtText,
		CloudSwiped: skyThought.CloudSwiped,
		SwipedAt:    skyThought.SwipedAt,
		CreatedDate: skyThought.CreatedDate,
		DayNumber:   skyThought.DayNumber,
		CreatedAt:   skyThought.CreatedAt,
		UpdatedAt:   skyThought.UpdatedAt,
	}
}
