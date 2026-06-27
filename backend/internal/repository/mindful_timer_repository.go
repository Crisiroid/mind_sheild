package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type MindfulTimerRepository struct {
	*BaseRepository
}

func NewMindfulTimerRepository(db *gorm.DB) *MindfulTimerRepository {
	return &MindfulTimerRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *MindfulTimerRepository) Create(ctx context.Context, timer *models.MindfulTimer) error {
	return r.BaseRepository.Create(ctx, timer)
}

func (r *MindfulTimerRepository) GetByID(ctx context.Context, id string) (*models.MindfulTimer, error) {
	var timer models.MindfulTimer
	err := r.BaseRepository.GetByID(ctx, &timer, id)
	if err != nil {
		return nil, err
	}
	return &timer, nil
}

func (r *MindfulTimerRepository) Update(ctx context.Context, timer *models.MindfulTimer) error {
	return r.BaseRepository.Update(ctx, timer)
}

func (r *MindfulTimerRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.MindfulTimer{}, id)
}

func (r *MindfulTimerRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MindfulTimer, int64, error) {
	var timers []models.MindfulTimer
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.MindfulTimer{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&timers).Error; err != nil {
		return nil, 0, err
	}

	return timers, total, nil
}

func (r *MindfulTimerRepository) GetByUserID(ctx context.Context, userID string) ([]models.MindfulTimer, error) {
	var timers []models.MindfulTimer
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&timers).Error
	return timers, err
}

func (r *MindfulTimerRepository) GetTimerStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*MindfulTimerStats, error) {
	var stats MindfulTimerStats

	query := r.DB.Model(&models.MindfulTimer{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalTimers)

	r.DB.Model(&models.MindfulTimer{}).
		Where("is_completed = ?", true).
		Count(&stats.CompletedTimers)

	stats.IncompleteTimers = stats.TotalTimers - stats.CompletedTimers
	if stats.TotalTimers > 0 {
		stats.CompletionRate = float64(stats.CompletedTimers) / float64(stats.TotalTimers) * 100
	}

	r.DB.Model(&models.MindfulTimer{}).
		Where("duration_seconds IS NOT NULL").
		Select("COALESCE(AVG(duration_seconds), 0)").
		Scan(&stats.AvgDuration)

	r.DB.Model(&models.MindfulTimer{}).
		Select("COALESCE(SUM(vibration_reminders_count), 0)").
		Scan(&stats.TotalVibrationReminders)

	return &stats, nil
}

func (r *MindfulTimerRepository) GetDurationTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count,
			AVG(duration_seconds) as value
		FROM mindful_timers
		WHERE user_id = ? AND duration_seconds IS NOT NULL
			AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *MindfulTimerRepository) GetVibrationAnalysis(ctx context.Context, userID string) (*schemas.TimeAnalysisResponse, error) {
	var stats schemas.TimeAnalysisResponse

	query := r.DB.Model(&models.MindfulTimer{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.EntriesCount)
	r.DB.Model(&models.MindfulTimer{}).
		Select("COALESCE(AVG(vibration_reminders_count), 0)").
		Scan(&stats.AvgDuration)
	r.DB.Model(&models.MindfulTimer{}).
		Select("COALESCE(MAX(vibration_reminders_count), 0)").
		Scan(&stats.MaxDuration)

	return &stats, nil
}

type MindfulTimerStats struct {
	TotalTimers             int64   `json:"total_timers"`
	CompletedTimers         int64   `json:"completed_timers"`
	IncompleteTimers        int64   `json:"incomplete_timers"`
	CompletionRate          float64 `json:"completion_rate"`
	AvgDuration             float64 `json:"avg_duration"`
	TotalVibrationReminders int64   `json:"total_vibration_reminders"`
}
