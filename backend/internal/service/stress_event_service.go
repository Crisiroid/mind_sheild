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

type StressEventService struct {
	stressRepo *repository.StressEventRepository
}

func NewStressEventService(stressRepo *repository.StressEventRepository) *StressEventService {
	return &StressEventService{
		stressRepo: stressRepo,
	}
}

func (s *StressEventService) CreateStressEvent(ctx context.Context, req *schemas.StressEventCreateRequest) (*schemas.StressEventResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.SituationType == "" {
		return nil, errors.New("نوع موقعیت الزامی است")
	}

	event := &models.StressEvent{
		UserID:               req.UserID,
		SituationType:        req.SituationType,
		SituationDescription: req.SituationDescription,
		IntensityLevel:       req.IntensityLevel,
		Location:             req.Location,
		DayNumber:            req.DayNumber,
	}

	if err := s.stressRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد رویداد استرس: %w", err)
	}

	return s.toStressEventResponse(event), nil
}

func (s *StressEventService) GetStressEventById(ctx context.Context, id string) (*schemas.StressEventResponse, error) {
	event, err := s.stressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("رویداد استرس مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت رویداد استرس: %w", err)
	}

	return s.toStressEventResponse(event), nil
}

func (s *StressEventService) GetAllStressEvents(ctx context.Context, req *schemas.StressEventListRequest) (*schemas.StressEventListResponse, error) {
	filterFunc := s.buildStressEventFilters(req)

	events, total, err := s.stressRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت رویدادهای استرس: %w", err)
	}

	responses := make([]schemas.StressEventResponse, len(events))
	for i, event := range events {
		responses[i] = *s.toStressEventResponse(&event)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.StressEventListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.StressEventResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *StressEventService) UpdateStressEvent(ctx context.Context, id string, req *schemas.StressEventUpdateRequest) (*schemas.StressEventResponse, error) {
	event, err := s.stressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("رویداد استرس مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت رویداد استرس: %w", err)
	}

	if req.IntensityLevel != nil {
		event.IntensityLevel = *req.IntensityLevel
	}

	if req.SituationDescription != nil {
		event.SituationDescription = *req.SituationDescription
	}

	if err := s.stressRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی رویداد استرس: %w", err)
	}

	return s.toStressEventResponse(event), nil
}

func (s *StressEventService) DeleteStressEvent(ctx context.Context, id string) error {
	_, err := s.stressRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("رویداد استرس مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت رویداد استرس: %w", err)
	}

	if err := s.stressRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف رویداد استرس: %w", err)
	}

	return nil
}

func (s *StressEventService) buildStressEventFilters(req *schemas.StressEventListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.SituationType != nil {
			db = db.Where("situation_type = ?", *req.SituationType)
		}

		if req.IntensityMin != nil {
			db = db.Where("intensity_level >= ?", *req.IntensityMin)
		}

		if req.IntensityMax != nil {
			db = db.Where("intensity_level <= ?", *req.IntensityMax)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *StressEventService) toStressEventResponse(event *models.StressEvent) *schemas.StressEventResponse {
	return &schemas.StressEventResponse{
		ID:                   event.ID.String(),
		UserID:               event.UserID,
		EventTimestamp:       event.EventTimestamp,
		SituationType:        event.SituationType,
		SituationDescription: event.SituationDescription,
		IntensityLevel:       event.IntensityLevel,
		Location:             event.Location,
		DayNumber:            event.DayNumber,
		CreatedAt:            event.CreatedAt,
		UpdatedAt:            event.UpdatedAt,
	}
}
