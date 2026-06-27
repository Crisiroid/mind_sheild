package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type UserReportRepository struct {
	*BaseRepository
}

func NewUserReportRepository(db *gorm.DB) *UserReportRepository {
	return &UserReportRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *UserReportRepository) Create(ctx context.Context, report *models.UserReport) error {
	return r.BaseRepository.Create(ctx, report)
}

func (r *UserReportRepository) GetByID(ctx context.Context, id string) (*models.UserReport, error) {
	var report models.UserReport
	err := r.BaseRepository.GetByID(ctx, &report, id)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *UserReportRepository) Update(ctx context.Context, report *models.UserReport) error {
	return r.BaseRepository.Update(ctx, report)
}

func (r *UserReportRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.UserReport{}, id)
}

func (r *UserReportRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.UserReport, int64, error) {
	var reports []models.UserReport
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.UserReport{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

func (r *UserReportRepository) GetReportStats(ctx context.Context, dateFrom, dateTo *time.Time) (*UserReportStats, error) {
	var stats UserReportStats

	query := r.DB.Model(&models.UserReport{})

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalReports)

	r.DB.Model(&models.UserReport{}).
		Select("COALESCE(SUM(total_users), 0)").
		Scan(&stats.TotalUsersTracked)

	r.DB.Model(&models.UserReport{}).
		Select("COALESCE(SUM(active_users), 0)").
		Scan(&stats.TotalActiveUsers)

	r.DB.Model(&models.UserReport{}).
		Select("COALESCE(AVG(avg_engagement_score), 0)").
		Scan(&stats.AvgEngagementScore)

	r.DB.Model(&models.UserReport{}).
		Select("COALESCE(SUM(crisis_alerts_count), 0)").
		Scan(&stats.TotalCrisisAlerts)

	return &stats, nil
}

func (r *UserReportRepository) GetReportTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM admin_panel.user_reports
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *UserReportRepository) GetReportsByType(ctx context.Context, reportType string) ([]models.UserReport, error) {
	var reports []models.UserReport
	err := r.DB.WithContext(ctx).
		Where("report_type = ?", reportType).
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}

func (r *UserReportRepository) GetEngagementTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			AVG(avg_engagement_score) as value
		FROM admin_panel.user_reports
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

type UserReportStats struct {
	TotalReports       int64   `json:"total_reports"`
	TotalUsersTracked  int64   `json:"total_users_tracked"`
	TotalActiveUsers   int64   `json:"total_active_users"`
	AvgEngagementScore float64 `json:"avg_engagement_score"`
	TotalCrisisAlerts  int64   `json:"total_crisis_alerts"`
}
