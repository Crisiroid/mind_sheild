package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type StressEventRepository struct {
	*BaseRepository
}

func NewStressEventRepository(db *gorm.DB) *StressEventRepository {
	return &StressEventRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *StressEventRepository) Create(ctx context.Context, event *models.StressEvent) error {
	return r.BaseRepository.Create(ctx, event)
}

func (r *StressEventRepository) GetByID(ctx context.Context, id string) (*models.StressEvent, error) {
	var event models.StressEvent
	err := r.BaseRepository.GetByID(ctx, &event, id)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *StressEventRepository) Update(ctx context.Context, event *models.StressEvent) error {
	return r.BaseRepository.Update(ctx, event)
}

func (r *StressEventRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.StressEvent{}, id)
}

func (r *StressEventRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.StressEvent, int64, error) {
	var events []models.StressEvent
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.StressEvent{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *StressEventRepository) GetByUserID(ctx context.Context, userID string) ([]models.StressEvent, error) {
	var events []models.StressEvent
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

func (r *StressEventRepository) GetStressStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*StressEventStats, error) {
	var stats StressEventStats

	query := r.DB.Model(&models.StressEvent{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("event_timestamp >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("event_timestamp <= ?", *dateTo)
	}

	query.Count(&stats.TotalEvents)

	r.DB.Model(&models.StressEvent{}).
		Select("COALESCE(AVG(intensity_level), 0)").
		Scan(&stats.AvgIntensity)

	r.DB.Model(&models.StressEvent{}).
		Where("intensity_level >= 8").
		Count(&stats.HighIntensityEvents)

	return &stats, nil
}

func (r *StressEventRepository) GetIntensityTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(event_timestamp) as timestamp,
			COUNT(*) as count,
			AVG(intensity_level) as value
		FROM stress_events
		WHERE user_id = ? AND event_timestamp >= ? AND event_timestamp <= ?
		GROUP BY DATE(event_timestamp)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *StressEventRepository) GetSituationTypeDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.StressEvent{}).
		Select("situation_type as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("situation_type").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *StressEventRepository) GetLocationHotspots(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.StressEvent{}).
		Select("location as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Where("location IS NOT NULL AND location != ''").
		Group("location").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *StressEventRepository) GetHighIntensityEvents(ctx context.Context, userID string, threshold int) ([]models.StressEvent, error) {
	var events []models.StressEvent
	query := r.DB.WithContext(ctx).
		Where("intensity_level >= ?", threshold)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("intensity_level DESC").Find(&events).Error
	return events, err
}

type StressEventStats struct {
	TotalEvents         int64   `json:"total_events"`
	AvgIntensity        float64 `json:"avg_intensity"`
	HighIntensityEvents int64   `json:"high_intensity_events"`
}
