package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type EmotionTriangleRepository struct {
	*BaseRepository
}

func NewEmotionTriangleRepository(db *gorm.DB) *EmotionTriangleRepository {
	return &EmotionTriangleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *EmotionTriangleRepository) Create(ctx context.Context, interaction *models.EmotionTriangleInteraction) error {
	return r.BaseRepository.Create(ctx, interaction)
}

func (r *EmotionTriangleRepository) GetByID(ctx context.Context, id string) (*models.EmotionTriangleInteraction, error) {
	var interaction models.EmotionTriangleInteraction
	err := r.BaseRepository.GetByID(ctx, &interaction, id)
	if err != nil {
		return nil, err
	}
	return &interaction, nil
}

func (r *EmotionTriangleRepository) Update(ctx context.Context, interaction *models.EmotionTriangleInteraction) error {
	return r.BaseRepository.Update(ctx, interaction)
}

func (r *EmotionTriangleRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.EmotionTriangleInteraction{}, id)
}

func (r *EmotionTriangleRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.EmotionTriangleInteraction, int64, error) {
	var interactions []models.EmotionTriangleInteraction
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.EmotionTriangleInteraction{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&interactions).Error; err != nil {
		return nil, 0, err
	}

	return interactions, total, nil
}

func (r *EmotionTriangleRepository) GetByUserID(ctx context.Context, userID string) ([]models.EmotionTriangleInteraction, error) {
	var interactions []models.EmotionTriangleInteraction
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&interactions).Error
	return interactions, err
}

func (r *EmotionTriangleRepository) GetInteractionStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*EmotionTriangleStats, error) {
	var stats EmotionTriangleStats

	query := r.DB.Model(&models.EmotionTriangleInteraction{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("interaction_date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("interaction_date <= ?", *dateTo)
	}

	query.Count(&stats.TotalInteractions)

	r.DB.Model(&models.EmotionTriangleInteraction{}).
		Where("vibration_triggered = ?", true).
		Count(&stats.VibrationTriggered)

	return &stats, nil
}

func (r *EmotionTriangleRepository) GetSideClickedDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.EmotionTriangleInteraction{}).
		Select("side_clicked as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("side_clicked").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *EmotionTriangleRepository) GetVibrationTriggerRate(ctx context.Context, userID string) (int64, int64, float64, error) {
	var total, triggered int64

	query := r.DB.Model(&models.EmotionTriangleInteraction{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	r.DB.Model(&models.EmotionTriangleInteraction{}).
		Where("vibration_triggered = ?", true).
		Count(&triggered)

	var rate float64
	if total > 0 {
		rate = float64(triggered) / float64(total) * 100
	}

	return total, triggered, rate, nil
}

func (r *EmotionTriangleRepository) GetInteractionTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(interaction_date) as timestamp,
			COUNT(*) as count
		FROM emotion_triangle_interactions
		WHERE user_id = ? AND interaction_date >= ? AND interaction_date <= ?
		GROUP BY DATE(interaction_date)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *EmotionTriangleRepository) GetThoughtAccountsAnalysis(ctx context.Context, userID string) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(interaction_date) as timestamp,
			COUNT(*) as count
		FROM emotion_triangle_interactions
		WHERE user_id = ?
			AND side_clicked = 'thought'
			AND thought_accounts_viewed IS NOT NULL
			AND thought_accounts_viewed != '[]'
		GROUP BY DATE(interaction_date)
		ORDER BY timestamp ASC
	`, userID).Scan(&trends)

	return trends, nil
}

type EmotionTriangleStats struct {
	TotalInteractions  int64 `json:"total_interactions"`
	VibrationTriggered int64 `json:"vibration_triggered"`
}
