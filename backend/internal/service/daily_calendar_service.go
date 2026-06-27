package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type DailyCalendarService struct {
	dailyCalendarRepo *repository.DailyCalendarRepository
}

func NewDailyCalendarService(dailyCalendarRepo *repository.DailyCalendarRepository) *DailyCalendarService {
	return &DailyCalendarService{
		dailyCalendarRepo: dailyCalendarRepo,
	}
}

func (s *DailyCalendarService) CreateCalendarEntry(ctx context.Context, req *schemas.CalendarCreateRequest) (*schemas.CalendarResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.DayNumber < 1 || req.DayNumber > 56 {
		return nil, errors.New("شماره روز باید بین 1 تا 56 باشد")
	}

	existing, err := s.dailyCalendarRepo.GetByUserAndDay(ctx, req.UserID, req.DayNumber)
	if err == nil && existing != nil {
		return nil, errors.New("این روز قبلاً ثبت شده است")
	}

	activitiesJSON := "[]"
	if len(req.ActivitiesCompleted) > 0 {
		if data, err := json.Marshal(req.ActivitiesCompleted); err == nil {
			activitiesJSON = string(data)
		}
	}

	calendar := &models.DailyCalendar{
		UserID:              req.UserID,
		DayNumber:           req.DayNumber,
		ActivitiesCompleted: activitiesJSON,
	}

	if err := s.dailyCalendarRepo.Create(ctx, calendar); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد تقویم روزانه: %w", err)
	}

	return s.toCalendarResponse(calendar), nil
}

func (s *DailyCalendarService) GetCalendarEntryById(ctx context.Context, id string) (*schemas.CalendarResponse, error) {
	calendar, err := s.dailyCalendarRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تقویم روزانه مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تقویم روزانه: %w", err)
	}

	return s.toCalendarResponse(calendar), nil
}

func (s *DailyCalendarService) GetAllCalendarEntries(ctx context.Context, req *schemas.CalendarListRequest) (*schemas.CalendarListResponse, error) {
	filterFunc := s.buildCalendarFilters(req)

	calendars, total, err := s.dailyCalendarRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تقویم‌های روزانه: %w", err)
	}

	responses := make([]schemas.CalendarResponse, len(calendars))
	for i, calendar := range calendars {
		responses[i] = *s.toCalendarResponse(&calendar)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.CalendarListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.CalendarResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *DailyCalendarService) UpdateCalendarEntry(ctx context.Context, id string, req *schemas.CalendarUpdateRequest) (*schemas.CalendarResponse, error) {
	calendar, err := s.dailyCalendarRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تقویم روزانه مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تقویم روزانه: %w", err)
	}

	if req.IsCompleted != nil {
		calendar.IsCompleted = *req.IsCompleted
	}

	if len(req.ActivitiesCompleted) > 0 {
		if data, err := json.Marshal(req.ActivitiesCompleted); err == nil {
			calendar.ActivitiesCompleted = string(data)
		}
	}

	if err := s.dailyCalendarRepo.Update(ctx, calendar); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی تقویم روزانه: %w", err)
	}

	return s.toCalendarResponse(calendar), nil
}

func (s *DailyCalendarService) DeleteCalendarEntry(ctx context.Context, id string) error {
	_, err := s.dailyCalendarRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("تقویم روزانه مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت تقویم روزانه: %w", err)
	}

	if err := s.dailyCalendarRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف تقویم روزانه: %w", err)
	}

	return nil
}

func (s *DailyCalendarService) GetCompletionStats(ctx context.Context, userID string) (*schemas.CompletionStatsResponse, error) {
	stats, err := s.dailyCalendarRepo.GetCompletionStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت آمار تکمیل: %w", err)
	}

	return stats, nil
}

func (s *DailyCalendarService) GetDayRangeProgress(ctx context.Context, userID string, fromDay, toDay int) (*schemas.CompletionStatsResponse, error) {
	progress, err := s.dailyCalendarRepo.GetDayRangeProgress(ctx, userID, fromDay, toDay)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت پیشرفت: %w", err)
	}

	return progress, nil
}

func (s *DailyCalendarService) GetStreakAnalysis(ctx context.Context, userID string) (*repository.StreakAnalysis, error) {
	analysis, err := s.dailyCalendarRepo.GetStreakAnalysis(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تحلیل پشت سر هم: %w", err)
	}

	return analysis, nil
}

func (s *DailyCalendarService) buildCalendarFilters(req *schemas.CalendarListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.IsCompleted != nil {
			db = db.Where("is_completed = ?", *req.IsCompleted)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("day_number ASC")
		}

		return db
	}
}

func (s *DailyCalendarService) toCalendarResponse(calendar *models.DailyCalendar) *schemas.CalendarResponse {
	var activities []string
	if calendar.ActivitiesCompleted != "" {
		json.Unmarshal([]byte(calendar.ActivitiesCompleted), &activities)
	}

	return &schemas.CalendarResponse{
		ID:                  calendar.ID.String(),
		UserID:              calendar.UserID,
		DayNumber:           calendar.DayNumber,
		CalendarDate:        calendar.CalendarDate.Format("2006-01-02"),
		IsCompleted:         calendar.IsCompleted,
		CompletedAt:         calendar.CompletedAt,
		ActivitiesCompleted: activities,
		CreatedAt:           calendar.CreatedAt,
		UpdatedAt:           calendar.UpdatedAt,
	}
}
