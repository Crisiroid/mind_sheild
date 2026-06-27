package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type SystemLogRepository struct {
	*BaseRepository
}

func NewSystemLogRepository(db *gorm.DB) *SystemLogRepository {
	return &SystemLogRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *SystemLogRepository) Create(ctx context.Context, log *models.SystemLog) error {
	return r.BaseRepository.Create(ctx, log)
}

func (r *SystemLogRepository) GetByID(ctx context.Context, id string) (*models.SystemLog, error) {
	var log models.SystemLog
	err := r.BaseRepository.GetByID(ctx, &log, id)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *SystemLogRepository) Update(ctx context.Context, log *models.SystemLog) error {
	return r.BaseRepository.Update(ctx, log)
}

func (r *SystemLogRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.SystemLog{}, id)
}

func (r *SystemLogRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.SystemLog{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *SystemLogRepository) GetLogStats(ctx context.Context, dateFrom, dateTo *time.Time) (*SystemLogStats, error) {
	var stats SystemLogStats

	query := r.DB.Model(&models.SystemLog{})

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalLogs)

	r.DB.Model(&models.SystemLog{}).
		Where("severity = ?", "info").
		Count(&stats.InfoCount)

	r.DB.Model(&models.SystemLog{}).
		Where("severity = ?", "warning").
		Count(&stats.WarningCount)

	r.DB.Model(&models.SystemLog{}).
		Where("severity = ?", "error").
		Count(&stats.ErrorCount)

	r.DB.Model(&models.SystemLog{}).
		Where("severity = ?", "critical").
		Count(&stats.CriticalCount)

	return &stats, nil
}

func (r *SystemLogRepository) GetSeverityDistribution(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	r.DB.Raw(`
		SELECT
			severity as label,
			COUNT(*) as count
		FROM admin_panel.system_logs
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY severity
		ORDER BY count DESC
	`, dateFrom, dateTo).Scan(&distributions)

	return distributions, nil
}

func (r *SystemLogRepository) GetLogTypeDistribution(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	r.DB.Raw(`
		SELECT
			log_type as label,
			COUNT(*) as count
		FROM admin_panel.system_logs
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY log_type
		ORDER BY count DESC
	`, dateFrom, dateTo).Scan(&distributions)

	return distributions, nil
}

func (r *SystemLogRepository) GetLogsBySeverity(ctx context.Context, severity string, limit int) ([]models.SystemLog, error) {
	var logs []models.SystemLog
	err := r.DB.WithContext(ctx).
		Where("severity = ?", severity).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

func (r *SystemLogRepository) GetErrorLogs(ctx context.Context, dateFrom, dateTo time.Time) ([]models.SystemLog, error) {
	var logs []models.SystemLog
	err := r.DB.WithContext(ctx).
		Where("severity IN ('error', 'critical')").
		Where("created_at >= ? AND created_at <= ?", dateFrom, dateTo).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *SystemLogRepository) GetLogTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM admin_panel.system_logs
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

type SystemLogStats struct {
	TotalLogs     int64 `json:"total_logs"`
	InfoCount     int64 `json:"info_count"`
	WarningCount  int64 `json:"warning_count"`
	ErrorCount    int64 `json:"error_count"`
	CriticalCount int64 `json:"critical_count"`
}
