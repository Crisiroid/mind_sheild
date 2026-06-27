package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type SkyThoughtRepository struct {
	*BaseRepository
}

func NewSkyThoughtRepository(db *gorm.DB) *SkyThoughtRepository {
	return &SkyThoughtRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *SkyThoughtRepository) Create(ctx context.Context, skyThought *models.SkyThought) error {
	return r.BaseRepository.Create(ctx, skyThought)
}

func (r *SkyThoughtRepository) GetByID(ctx context.Context, id string) (*models.SkyThought, error) {
	var skyThought models.SkyThought
	err := r.BaseRepository.GetByID(ctx, &skyThought, id)
	if err != nil {
		return nil, err
	}
	return &skyThought, nil
}

func (r *SkyThoughtRepository) Update(ctx context.Context, skyThought *models.SkyThought) error {
	return r.BaseRepository.Update(ctx, skyThought)
}

func (r *SkyThoughtRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.SkyThought{}, id)
}

func (r *SkyThoughtRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.SkyThought, int64, error) {
	var skyThoughts []models.SkyThought
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.SkyThought{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&skyThoughts).Error; err != nil {
		return nil, 0, err
	}

	return skyThoughts, total, nil
}

func (r *SkyThoughtRepository) GetByUserID(ctx context.Context, userID string) ([]models.SkyThought, error) {
	var skyThoughts []models.SkyThought
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&skyThoughts).Error
	return skyThoughts, err
}

func (r *SkyThoughtRepository) GetSkyThoughtStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*SkyThoughtStats, error) {
	var stats SkyThoughtStats

	query := r.DB.Model(&models.SkyThought{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalThoughts)

	r.DB.Model(&models.SkyThought{}).
		Where("cloud_swiped = ?", true).
		Count(&stats.CloudSwiped)

	if stats.TotalThoughts > 0 {
		stats.SwipeRate = float64(stats.CloudSwiped) / float64(stats.TotalThoughts) * 100
	}

	return &stats, nil
}

func (r *SkyThoughtRepository) GetSwipeRate(ctx context.Context, userID string) (int64, int64, float64, error) {
	var total, swiped int64

	query := r.DB.Model(&models.SkyThought{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	r.DB.Model(&models.SkyThought{}).
		Where("cloud_swiped = ?", true).
		Count(&swiped)

	var rate float64
	if total > 0 {
		rate = float64(swiped) / float64(total) * 100
	}

	return total, swiped, rate, nil
}

func (r *SkyThoughtRepository) GetThoughtTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM sky_thoughts
		WHERE user_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *SkyThoughtRepository) GetUnswipedThoughts(ctx context.Context, userID string) ([]models.SkyThought, error) {
	var skyThoughts []models.SkyThought
	query := r.DB.WithContext(ctx).Where("cloud_swiped = ?", false)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("created_at ASC").Find(&skyThoughts).Error
	return skyThoughts, err
}

type SkyThoughtStats struct {
	TotalThoughts int64   `json:"total_thoughts"`
	CloudSwiped   int64   `json:"cloud_swiped"`
	SwipeRate     float64 `json:"swipe_rate"`
}
