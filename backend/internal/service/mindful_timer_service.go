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

type MindfulTimerService struct {
	mindfulTimerRepo *repository.MindfulTimerRepository
}

func NewMindfulTimerService(mindfulTimerRepo *repository.MindfulTimerRepository) *MindfulTimerService {
	return &MindfulTimerService{
		mindfulTimerRepo: mindfulTimerRepo,
	}
}

func (s *MindfulTimerService) CreateMindfulTimer(ctx context.Context, req *schemas.MindfulTimerCreateRequest) (*schemas.MindfulTimerResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	timer := &models.MindfulTimer{
		UserID:    req.UserID,
		DayNumber: req.DayNumber,
	}

	if err := s.mindfulTimerRepo.Create(ctx, timer); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد تایمر ذهن‌آگاهی: %w", err)
	}

	return s.toMindfulTimerResponse(timer), nil
}

func (s *MindfulTimerService) GetMindfulTimerById(ctx context.Context, id string) (*schemas.MindfulTimerResponse, error) {
	timer, err := s.mindfulTimerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تایمر ذهن‌آگاهی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تایمر ذهن‌آگاهی: %w", err)
	}

	return s.toMindfulTimerResponse(timer), nil
}

func (s *MindfulTimerService) GetAllMindfulTimers(ctx context.Context, req *schemas.MindfulTimerListRequest) (*schemas.MindfulTimerListResponse, error) {
	filterFunc := s.buildMindfulTimerFilters(req)

	timers, total, err := s.mindfulTimerRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تایمر‌های ذهن‌آگاهی: %w", err)
	}

	responses := make([]schemas.MindfulTimerResponse, len(timers))
	for i, timer := range timers {
		responses[i] = *s.toMindfulTimerResponse(&timer)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.MindfulTimerListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.MindfulTimerResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *MindfulTimerService) UpdateMindfulTimer(ctx context.Context, id string, req *schemas.MindfulTimerUpdateRequest) (*schemas.MindfulTimerResponse, error) {
	timer, err := s.mindfulTimerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تایمر ذهن‌آگاهی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تایمر ذهن‌آگاهی: %w", err)
	}

	if req.TimerEnd != nil {
		timer.TimerEnd = *req.TimerEnd
	}

	if req.DurationSeconds != nil {
		timer.DurationSeconds = req.DurationSeconds
	}

	if req.VibrationRemindersCount != nil {
		timer.VibrationRemindersCount = *req.VibrationRemindersCount
	}

	if req.IsCompleted != nil {
		timer.IsCompleted = *req.IsCompleted
	}

	if err := s.mindfulTimerRepo.Update(ctx, timer); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی تایمر ذهن‌آگاهی: %w", err)
	}

	return s.toMindfulTimerResponse(timer), nil
}

func (s *MindfulTimerService) DeleteMindfulTimer(ctx context.Context, id string) error {
	_, err := s.mindfulTimerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("تایمر ذهن‌آگاهی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت تایمر ذهن‌آگاهی: %w", err)
	}

	if err := s.mindfulTimerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف تایمر ذهن‌آگاهی: %w", err)
	}

	return nil
}

func (s *MindfulTimerService) buildMindfulTimerFilters(req *schemas.MindfulTimerListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.IsCompleted != nil {
			db = db.Where("is_completed = ?", *req.IsCompleted)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *MindfulTimerService) toMindfulTimerResponse(timer *models.MindfulTimer) *schemas.MindfulTimerResponse {
	return &schemas.MindfulTimerResponse{
		ID:                      timer.ID.String(),
		UserID:                  timer.UserID,
		TimerStart:              timer.TimerStart,
		TimerEnd:                timer.TimerEnd,
		DurationSeconds:         timer.DurationSeconds,
		VibrationRemindersCount: timer.VibrationRemindersCount,
		IsCompleted:             timer.IsCompleted,
		DayNumber:               timer.DayNumber,
		CreatedAt:               timer.CreatedAt,
		UpdatedAt:               timer.UpdatedAt,
	}
}
