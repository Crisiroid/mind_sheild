package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type MentalMustRepository struct {
	*BaseRepository
}

func NewMentalMustRepository(db *gorm.DB) *MentalMustRepository {
	return &MentalMustRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *MentalMustRepository) Create(ctx context.Context, mentalMust *models.MentalMust) error {
	return r.BaseRepository.Create(ctx, mentalMust)
}

func (r *MentalMustRepository) GetByID(ctx context.Context, id string) (*models.MentalMust, error) {
	var mentalMust models.MentalMust
	err := r.BaseRepository.GetByID(ctx, &mentalMust, id)
	if err != nil {
		return nil, err
	}
	return &mentalMust, nil
}

func (r *MentalMustRepository) Update(ctx context.Context, mentalMust *models.MentalMust) error {
	return r.BaseRepository.Update(ctx, mentalMust)
}

func (r *MentalMustRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.MentalMust{}, id)
}

func (r *MentalMustRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MentalMust, int64, error) {
	var mentalMusts []models.MentalMust
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.MentalMust{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&mentalMusts).Error; err != nil {
		return nil, 0, err
	}

	return mentalMusts, total, nil
}

func (r *MentalMustRepository) GetByUserID(ctx context.Context, userID string) ([]models.MentalMust, error) {
	var mentalMusts []models.MentalMust
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&mentalMusts).Error
	return mentalMusts, err
}

func (r *MentalMustRepository) GetMustStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*MentalMustStats, error) {
	var stats MentalMustStats

	query := r.DB.Model(&models.MentalMust{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalMusts)

	r.DB.Model(&models.MentalMust{}).
		Where("is_released = ?", true).
		Count(&stats.ReleasedMusts)

	stats.UnreleasedMusts = stats.TotalMusts - stats.ReleasedMusts
	if stats.TotalMusts > 0 {
		stats.ReleaseRate = float64(stats.ReleasedMusts) / float64(stats.TotalMusts) * 100
	}

	return &stats, nil
}

func (r *MentalMustRepository) GetReleaseRate(ctx context.Context, userID string) (int64, int64, float64, error) {
	var total, released int64

	query := r.DB.Model(&models.MentalMust{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	r.DB.Model(&models.MentalMust{}).
		Where("is_released = ?", true).
		Count(&released)

	var rate float64
	if total > 0 {
		rate = float64(released) / float64(total) * 100
	}

	return total, released, rate, nil
}

func (r *MentalMustRepository) GetReleaseTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM mental_musts
		WHERE is_released = true
			AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *MentalMustRepository) GetUserMustCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.MentalMust{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *MentalMustRepository) GetMustsByReleaseStatus(ctx context.Context, userID string, isReleased bool) ([]models.MentalMust, error) {
	var mentalMusts []models.MentalMust
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND is_released = ?", userID, isReleased).
		Order("created_at DESC").
		Find(&mentalMusts).Error
	return mentalMusts, err
}

type MentalMustStats struct {
	TotalMusts      int64   `json:"total_musts"`
	ReleasedMusts   int64   `json:"released_musts"`
	UnreleasedMusts int64   `json:"unreleased_musts"`
	ReleaseRate     float64 `json:"release_rate"`
}
