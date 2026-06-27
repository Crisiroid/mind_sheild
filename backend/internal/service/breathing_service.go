package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type BreathingService struct {
	breathingRepo interfaces.BreathingRepositoryInterface
}

func NewBreathingService(breathingRepo interfaces.BreathingRepositoryInterface) *BreathingService {
	return &BreathingService{
		breathingRepo: breathingRepo,
	}
}

func (s *BreathingService) CreateBreathingSession(ctx context.Context, req *schemas.BreathingSessionCreateRequest) (*schemas.BreathingSessionResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	session := &models.BreathingSession{
		UserID:           req.UserID,
		BreathingPattern: req.BreathingPattern,
		DayNumber:        req.DayNumber,
	}

	if err := s.breathingRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد جلسه تنفس: %w", err)
	}

	return s.toBreathingSessionResponse(session), nil
}

func (s *BreathingService) GetBreathingSessionById(ctx context.Context, id string) (*schemas.BreathingSessionResponse, error) {
	session, err := s.breathingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("جلسه تنفس مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت جلسه تنفس: %w", err)
	}

	return s.toBreathingSessionResponse(session), nil
}

func (s *BreathingService) GetAllBreathingSessions(ctx context.Context, req *schemas.BreathingSessionListRequest) (*schemas.BreathingSessionListResponse, error) {
	filterFunc := s.buildBreathingSessionFilters(req)

	sessions, total, err := s.breathingRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت جلسات تنفس: %w", err)
	}

	responses := make([]schemas.BreathingSessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = *s.toBreathingSessionResponse(&session)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.BreathingSessionListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.BreathingSessionResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *BreathingService) UpdateBreathingSession(ctx context.Context, id string, req *schemas.BreathingSessionUpdateRequest) (*schemas.BreathingSessionResponse, error) {
	session, err := s.breathingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("جلسه تنفس مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت جلسه تنفس: %w", err)
	}

	if req.SessionEnd != nil {
		session.SessionEnd = req.SessionEnd
	}

	if req.DurationSeconds != nil {
		session.DurationSeconds = req.DurationSeconds
	}

	if req.IsCompleted != nil {
		session.IsCompleted = *req.IsCompleted
	}

	if req.CalendarTicked != nil {
		session.CalendarTicked = *req.CalendarTicked
	}

	if err := s.breathingRepo.Update(ctx, session); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی جلسه تنفس: %w", err)
	}

	return s.toBreathingSessionResponse(session), nil
}

func (s *BreathingService) DeleteBreathingSession(ctx context.Context, id string) error {
	_, err := s.breathingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("جلسه تنفس مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت جلسه تنفس: %w", err)
	}

	if err := s.breathingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف جلسه تنفس: %w", err)
	}

	return nil
}

func (s *BreathingService) buildBreathingSessionFilters(req *schemas.BreathingSessionListRequest) func(*gorm.DB) *gorm.DB {
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

func (s *BreathingService) GetSessionStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.BreathingSessionStats, error) {
	return s.breathingRepo.GetSessionStats(ctx, userID, dateFrom, dateTo)
}

func (s *BreathingService) GetDurationTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	return s.breathingRepo.GetDurationTrend(ctx, userID, dateFrom, dateTo)
}

func (s *BreathingService) GetPatternUsage(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	return s.breathingRepo.GetPatternUsage(ctx, userID)
}

func (s *BreathingService) toBreathingSessionResponse(session *models.BreathingSession) *schemas.BreathingSessionResponse {
	return &schemas.BreathingSessionResponse{
		ID:               session.ID.String(),
		UserID:           session.UserID,
		SessionStart:     session.SessionStart,
		SessionEnd:       session.SessionEnd,
		DurationSeconds:  session.DurationSeconds,
		BreathingPattern: session.BreathingPattern,
		IsCompleted:      session.IsCompleted,
		CalendarTicked:   session.CalendarTicked,
		DayNumber:        session.DayNumber,
		CreatedAt:        session.CreatedAt,
		UpdatedAt:        session.UpdatedAt,
	}
}
