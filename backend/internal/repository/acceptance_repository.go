package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type AcceptanceRepository struct {
	*BaseRepository
}

func NewAcceptanceRepository(db *gorm.DB) *AcceptanceRepository {
	return &AcceptanceRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *AcceptanceRepository) Create(ctx context.Context, exercise *models.AcceptanceExercise) error {
	return r.BaseRepository.Create(ctx, exercise)
}

func (r *AcceptanceRepository) GetByID(ctx context.Context, id string) (*models.AcceptanceExercise, error) {
	var exercise models.AcceptanceExercise
	err := r.BaseRepository.GetByID(ctx, &exercise, id)
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *AcceptanceRepository) Update(ctx context.Context, exercise *models.AcceptanceExercise) error {
	return r.BaseRepository.Update(ctx, exercise)
}

func (r *AcceptanceRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.AcceptanceExercise{}, id)
}

func (r *AcceptanceRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.AcceptanceExercise, int64, error) {
	var exercises []models.AcceptanceExercise
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.AcceptanceExercise{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&exercises).Error; err != nil {
		return nil, 0, err
	}

	return exercises, total, nil
}

func (r *AcceptanceRepository) GetByUserID(ctx context.Context, userID string) ([]models.AcceptanceExercise, error) {
	var exercises []models.AcceptanceExercise
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&exercises).Error
	return exercises, err
}

func (r *AcceptanceRepository) GetCompletionStats(ctx context.Context, dateFrom, dateTo *time.Time, userID string) (*schemas.CompletionStatsResponse, error) {
	var stats schemas.CompletionStatsResponse

	query := r.DB.Model(&models.AcceptanceExercise{})

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.Total)

	r.DB.Model(&models.AcceptanceExercise{}).
		Where("video_watched = ?", true).
		Count(&stats.Completed)

	stats.Incomplete = stats.Total - stats.Completed
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	return &stats, nil
}

func (r *AcceptanceRepository) GetUnderstandingLevelDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.AcceptanceExercise{}).
		Select("understanding_level as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("understanding_level").
		Order("understanding_level ASC").
		Scan(&distributions)

	return distributions, nil
}

func (r *AcceptanceRepository) GetExerciseProgress(ctx context.Context, userID string) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			created_at as timestamp,
			COUNT(*) as count,
			AVG(understanding_level) as value
		FROM acceptance_exercises
		WHERE user_id = ?
		GROUP BY created_at
		ORDER BY timestamp ASC
	`, userID).Scan(&trends)

	return trends, nil
}

func (r *AcceptanceRepository) GetNotesAnalytics(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM acceptance_exercises
		WHERE notes IS NOT NULL AND notes != ''
			AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *AcceptanceRepository) GetAverageUnderstandingLevel(ctx context.Context, userID string) (float64, error) {
	var avg float64

	query := r.DB.Model(&models.AcceptanceExercise{}).
		Select("COALESCE(AVG(understanding_level), 0)")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Scan(&avg)
	return avg, nil
}

func (r *AcceptanceRepository) GetVideoWatchedStats(ctx context.Context, userID string) (int64, int64, float64, error) {
	var total, watched int64

	query := r.DB.Model(&models.AcceptanceExercise{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	query.Where("video_watched = ?", true).Count(&watched)

	var rate float64
	if total > 0 {
		rate = float64(watched) / float64(total) * 100
	}

	return total, watched, rate, nil
}
