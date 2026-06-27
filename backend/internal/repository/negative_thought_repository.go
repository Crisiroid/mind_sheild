package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type NegativeThoughtRepository struct {
	*BaseRepository
}

func NewNegativeThoughtRepository(db *gorm.DB) *NegativeThoughtRepository {
	return &NegativeThoughtRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *NegativeThoughtRepository) Create(ctx context.Context, thought *models.NegativeThought) error {
	return r.BaseRepository.Create(ctx, thought)
}

func (r *NegativeThoughtRepository) GetByID(ctx context.Context, id string) (*models.NegativeThought, error) {
	var thought models.NegativeThought
	err := r.BaseRepository.GetByID(ctx, &thought, id)
	if err != nil {
		return nil, err
	}
	return &thought, nil
}

func (r *NegativeThoughtRepository) Update(ctx context.Context, thought *models.NegativeThought) error {
	return r.BaseRepository.Update(ctx, thought)
}

func (r *NegativeThoughtRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.NegativeThought{}, id)
}

func (r *NegativeThoughtRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.NegativeThought, int64, error) {
	var thoughts []models.NegativeThought
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.NegativeThought{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&thoughts).Error; err != nil {
		return nil, 0, err
	}

	return thoughts, total, nil
}

func (r *NegativeThoughtRepository) GetByUserID(ctx context.Context, userID string) ([]models.NegativeThought, error) {
	var thoughts []models.NegativeThought
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&thoughts).Error
	return thoughts, err
}

func (r *NegativeThoughtRepository) GetThoughtStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*NegativeThoughtStats, error) {
	var stats NegativeThoughtStats

	query := r.DB.Model(&models.NegativeThought{})

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

	r.DB.Model(&models.NegativeThought{}).
		Where("impact_level IS NOT NULL").
		Select("COALESCE(AVG(impact_level), 0)").
		Scan(&stats.AvgImpactLevel)

	return &stats, nil
}

func (r *NegativeThoughtRepository) GetCognitiveErrorDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.NegativeThought{}).
		Select("cognitive_error_type as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("cognitive_error_type").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *NegativeThoughtRepository) GetImpactTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count,
			AVG(impact_level) as value
		FROM negative_thoughts
		WHERE user_id = ? AND impact_level IS NOT NULL
			AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *NegativeThoughtRepository) GetHighImpactThoughts(ctx context.Context, userID string, threshold int) ([]models.NegativeThought, error) {
	var thoughts []models.NegativeThought
	query := r.DB.WithContext(ctx).
		Where("impact_level >= ?", threshold)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("impact_level DESC").Find(&thoughts).Error
	return thoughts, err
}

type NegativeThoughtStats struct {
	TotalThoughts  int64   `json:"total_thoughts"`
	AvgImpactLevel float64 `json:"avg_impact_level"`
}
