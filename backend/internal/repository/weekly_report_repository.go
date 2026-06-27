package repository

import (
	"context"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type WeeklyReportRepository struct {
	*BaseRepository
}

func NewWeeklyReportRepository(db *gorm.DB) *WeeklyReportRepository {
	return &WeeklyReportRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WeeklyReportRepository) Create(ctx context.Context, report *models.WeeklyReport) error {
	return r.BaseRepository.Create(ctx, report)
}

func (r *WeeklyReportRepository) GetByID(ctx context.Context, id string) (*models.WeeklyReport, error) {
	var report models.WeeklyReport
	err := r.BaseRepository.GetByID(ctx, &report, id)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *WeeklyReportRepository) GetByUserAndWeek(ctx context.Context, userID string, weekNumber int) (*models.WeeklyReport, error) {
	var report models.WeeklyReport
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND week_number = ?", userID, weekNumber).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *WeeklyReportRepository) Update(ctx context.Context, report *models.WeeklyReport) error {
	return r.BaseRepository.Update(ctx, report)
}

func (r *WeeklyReportRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.WeeklyReport{}, id)
}

func (r *WeeklyReportRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.WeeklyReport, int64, error) {
	var reports []models.WeeklyReport
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.WeeklyReport{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("week_number ASC").Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *WeeklyReportRepository) GetByUserID(ctx context.Context, userID string) ([]models.WeeklyReport, error) {
	var reports []models.WeeklyReport
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("week_number ASC").
		Find(&reports).Error
	return reports, err
}

func (r *WeeklyReportRepository) GetWeeklyStats(ctx context.Context, userID string) (*WeeklyReportStats, error) {
	var stats WeeklyReportStats

	query := r.DB.Model(&models.WeeklyReport{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.TotalReports)

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(AVG(progress_percentage), 0)").
		Scan(&stats.AvgProgress)

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(SUM(stress_events_count), 0)").
		Scan(&stats.TotalStressEvents)

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(SUM(breathing_sessions_count), 0)").
		Scan(&stats.TotalBreathingSessions)

	return &stats, nil
}

func (r *WeeklyReportRepository) GetProgressTrend(ctx context.Context, userID string) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			week_number as count,
			progress_percentage as value
		FROM weekly_reports
		WHERE user_id = ?
		ORDER BY week_number ASC
	`, userID).Scan(&trends)

	return trends, nil
}

func (r *WeeklyReportRepository) GetActivityDistribution(ctx context.Context, userID string) (*WeeklyActivityDistribution, error) {
	var dist WeeklyActivityDistribution

	query := r.DB.Model(&models.WeeklyReport{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(SUM(stress_events_count), 0)").
		Scan(&dist.TotalStressEvents)

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(SUM(breathing_sessions_count), 0)").
		Scan(&dist.TotalBreathingSessions)

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(SUM(negative_thoughts_count), 0)").
		Scan(&dist.TotalNegativeThoughts)

	r.DB.Model(&models.WeeklyReport{}).
		Select("COALESCE(SUM(body_tension_maps_count), 0)").
		Scan(&dist.TotalBodyTensionMaps)

	return &dist, nil
}

func (r *WeeklyReportRepository) GetWeeklyComparison(ctx context.Context, userID string, week1, week2 int) (*WeeklyComparison, error) {
	var report1, report2 models.WeeklyReport

	r.DB.WithContext(ctx).
		Where("user_id = ? AND week_number = ?", userID, week1).
		First(&report1)

	r.DB.WithContext(ctx).
		Where("user_id = ? AND week_number = ?", userID, week2).
		First(&report2)

	comparison := &WeeklyComparison{
		Week1:            week1,
		Week2:            week2,
		StressEventsDiff: report2.StressEventsCount - report1.StressEventsCount,
		BreathingDiff:    report2.BreathingSessionsCount - report1.BreathingSessionsCount,
		ThoughtsDiff:     report2.NegativeThoughtsCount - report1.NegativeThoughtsCount,
		ProgressDiff:     report2.ProgressPercentage - report1.ProgressPercentage,
	}

	return comparison, nil
}

type WeeklyReportStats struct {
	TotalReports           int64   `json:"total_reports"`
	AvgProgress            float64 `json:"avg_progress"`
	TotalStressEvents      int64   `json:"total_stress_events"`
	TotalBreathingSessions int64   `json:"total_breathing_sessions"`
}

type WeeklyActivityDistribution struct {
	TotalStressEvents      int64 `json:"total_stress_events"`
	TotalBreathingSessions int64 `json:"total_breathing_sessions"`
	TotalNegativeThoughts  int64 `json:"total_negative_thoughts"`
	TotalBodyTensionMaps   int64 `json:"total_body_tension_maps"`
}

type WeeklyComparison struct {
	Week1            int     `json:"week_1"`
	Week2            int     `json:"week_2"`
	StressEventsDiff int     `json:"stress_events_diff"`
	BreathingDiff    int     `json:"breathing_diff"`
	ThoughtsDiff     int     `json:"thoughts_diff"`
	ProgressDiff     float64 `json:"progress_diff"`
}
