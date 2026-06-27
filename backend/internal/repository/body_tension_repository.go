package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type BodyTensionRepository struct {
	*BaseRepository
}

func NewBodyTensionRepository(db *gorm.DB) *BodyTensionRepository {
	return &BodyTensionRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *BodyTensionRepository) Create(ctx context.Context, tensionMap *models.BodyTensionMap) error {
	return r.BaseRepository.Create(ctx, tensionMap)
}

func (r *BodyTensionRepository) GetByID(ctx context.Context, id string) (*models.BodyTensionMap, error) {
	var tensionMap models.BodyTensionMap
	err := r.BaseRepository.GetByID(ctx, &tensionMap, id)
	if err != nil {
		return nil, err
	}
	return &tensionMap, nil
}

func (r *BodyTensionRepository) Update(ctx context.Context, tensionMap *models.BodyTensionMap) error {
	return r.BaseRepository.Update(ctx, tensionMap)
}

func (r *BodyTensionRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.BodyTensionMap{}, id)
}

func (r *BodyTensionRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.BodyTensionMap, int64, error) {
	var tensionMaps []models.BodyTensionMap
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.BodyTensionMap{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&tensionMaps).Error; err != nil {
		return nil, 0, err
	}

	return tensionMaps, total, nil
}

func (r *BodyTensionRepository) GetByUserID(ctx context.Context, userID string) ([]models.BodyTensionMap, error) {
	var tensionMaps []models.BodyTensionMap
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tensionMaps).Error
	return tensionMaps, err
}

func (r *BodyTensionRepository) GetIntensityTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	query := r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count,
			AVG(overall_intensity) as value
		FROM body_tension_maps
		WHERE user_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo)

	query.Scan(&trends)
	return trends, nil
}

func (r *BodyTensionRepository) GetSeverityColorDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.BodyTensionMap{}).
		Select("severity_color as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("severity_color").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *BodyTensionRepository) GetBodyRegionHeatmap(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.BodyTensionMap{}).
		Select("id as label, overall_intensity as value")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Order("overall_intensity DESC").
		Limit(20).
		Scan(&distributions)

	return distributions, nil
}

func (r *BodyTensionRepository) GetAverageIntensity(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*schemas.IntensityStatsResponse, error) {
	var stats schemas.IntensityStatsResponse

	query := r.DB.Model(&models.BodyTensionMap{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalEntries)

	r.DB.Model(&models.BodyTensionMap{}).
		Select("COALESCE(AVG(overall_intensity), 0)").
		Scan(&stats.AvgIntensity)

	r.DB.Model(&models.BodyTensionMap{}).
		Select("COALESCE(MIN(overall_intensity), 0)").
		Scan(&stats.MinIntensity)

	r.DB.Model(&models.BodyTensionMap{}).
		Select("COALESCE(MAX(overall_intensity), 0)").
		Scan(&stats.MaxIntensity)

	r.DB.Model(&models.BodyTensionMap{}).
		Where("overall_intensity >= 8").
		Count(&stats.HighIntensity)

	r.DB.Model(&models.BodyTensionMap{}).
		Where("overall_intensity >= 4 AND overall_intensity < 8").
		Count(&stats.MediumIntensity)

	r.DB.Model(&models.BodyTensionMap{}).
		Where("overall_intensity < 4").
		Count(&stats.LowIntensity)

	return &stats, nil
}

func (r *BodyTensionRepository) GetTensionMapsByDateRange(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]models.BodyTensionMap, error) {
	var tensionMaps []models.BodyTensionMap
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, dateFrom, dateTo).
		Order("created_at DESC").
		Find(&tensionMaps).Error
	return tensionMaps, err
}
