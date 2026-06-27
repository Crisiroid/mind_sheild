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

type WeeklyReportService struct {
	weeklyReportRepo interfaces.WeeklyReportRepositoryInterface
}

func NewWeeklyReportService(weeklyReportRepo interfaces.WeeklyReportRepositoryInterface) *WeeklyReportService {
	return &WeeklyReportService{
		weeklyReportRepo: weeklyReportRepo,
	}
}

func (s *WeeklyReportService) CreateWeeklyReport(ctx context.Context, req *schemas.WeeklyReportCreateRequest) (*schemas.WeeklyReportResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.WeekNumber < 1 || req.WeekNumber > 8 {
		return nil, errors.New("شماره هفته باید بین 1 تا 8 باشد")
	}

	existing, err := s.weeklyReportRepo.GetByUserAndWeek(ctx, req.UserID, req.WeekNumber)
	if err == nil && existing != nil {
		return nil, errors.New("گزارش این هفته قبلاً ثبت شده است")
	}

	report := &models.WeeklyReport{
		UserID:                 req.UserID,
		WeekNumber:             req.WeekNumber,
		StartDate:              req.StartDate,
		EndDate:                req.EndDate,
		StressEventsCount:      req.StressEventsCount,
		AvgStressIntensity:     req.AvgStressIntensity,
		BreathingSessionsCount: req.BreathingSessionsCount,
		NegativeThoughtsCount:  req.NegativeThoughtsCount,
		BodyTensionMapsCount:   req.BodyTensionMapsCount,
		MoodImprovementScore:   req.MoodImprovementScore,
		ActivitiesDistribution: req.ActivitiesDistribution,
		ProgressPercentage:     req.ProgressPercentage,
	}

	if err := s.weeklyReportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد گزارش هفتگی: %w", err)
	}

	return s.toWeeklyReportResponse(report), nil
}

func (s *WeeklyReportService) GetWeeklyReportById(ctx context.Context, id string) (*schemas.WeeklyReportResponse, error) {
	report, err := s.weeklyReportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("گزارش هفتگی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت گزارش هفتگی: %w", err)
	}

	return s.toWeeklyReportResponse(report), nil
}

func (s *WeeklyReportService) GetAllWeeklyReports(ctx context.Context, req *schemas.WeeklyReportListRequest) (*schemas.WeeklyReportListResponse, error) {
	filterFunc := s.buildWeeklyReportFilters(req)

	reports, total, err := s.weeklyReportRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت گزارش‌های هفتگی: %w", err)
	}

	responses := make([]schemas.WeeklyReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = *s.toWeeklyReportResponse(&report)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.WeeklyReportListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.WeeklyReportResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *WeeklyReportService) UpdateWeeklyReport(ctx context.Context, id string, req *schemas.WeeklyReportUpdateRequest) (*schemas.WeeklyReportResponse, error) {
	report, err := s.weeklyReportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("گزارش هفتگی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت گزارش هفتگی: %w", err)
	}

	if req.StressEventsCount != nil {
		report.StressEventsCount = *req.StressEventsCount
	}

	if req.AvgStressIntensity != nil {
		report.AvgStressIntensity = *req.AvgStressIntensity
	}

	if req.BreathingSessionsCount != nil {
		report.BreathingSessionsCount = *req.BreathingSessionsCount
	}

	if req.NegativeThoughtsCount != nil {
		report.NegativeThoughtsCount = *req.NegativeThoughtsCount
	}

	if req.BodyTensionMapsCount != nil {
		report.BodyTensionMapsCount = *req.BodyTensionMapsCount
	}

	if req.MoodImprovementScore != nil {
		report.MoodImprovementScore = *req.MoodImprovementScore
	}

	if req.ProgressPercentage != nil {
		report.ProgressPercentage = *req.ProgressPercentage
	}

	if err := s.weeklyReportRepo.Update(ctx, report); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی گزارش هفتگی: %w", err)
	}

	return s.toWeeklyReportResponse(report), nil
}

func (s *WeeklyReportService) DeleteWeeklyReport(ctx context.Context, id string) error {
	_, err := s.weeklyReportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("گزارش هفتگی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت گزارش هفتگی: %w", err)
	}

	if err := s.weeklyReportRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف گزارش هفتگی: %w", err)
	}

	return nil
}

func (s *WeeklyReportService) GetWeeklyStats(ctx context.Context, userID string) (*repository.WeeklyReportStats, error) {
	stats, err := s.weeklyReportRepo.GetWeeklyStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت آمار هفتگی: %w", err)
	}

	return stats, nil
}

func (s *WeeklyReportService) buildWeeklyReportFilters(req *schemas.WeeklyReportListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.WeekNumber != nil {
			db = db.Where("week_number = ?", *req.WeekNumber)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("week_number ASC")
		}

		return db
	}
}

func (s *WeeklyReportService) toWeeklyReportResponse(report *models.WeeklyReport) *schemas.WeeklyReportResponse {
	return &schemas.WeeklyReportResponse{
		ID:                     report.ID.String(),
		UserID:                 report.UserID,
		WeekNumber:             report.WeekNumber,
		StartDate:              report.StartDate,
		EndDate:                report.EndDate,
		StressEventsCount:      report.StressEventsCount,
		AvgStressIntensity:     report.AvgStressIntensity,
		BreathingSessionsCount: report.BreathingSessionsCount,
		NegativeThoughtsCount:  report.NegativeThoughtsCount,
		BodyTensionMapsCount:   report.BodyTensionMapsCount,
		MoodImprovementScore:   report.MoodImprovementScore,
		ActivitiesDistribution: report.ActivitiesDistribution,
		ProgressPercentage:     report.ProgressPercentage,
		GeneratedAt:            report.GeneratedAt,
		CreatedAt:              report.CreatedAt,
		UpdatedAt:              report.UpdatedAt,
	}
}
