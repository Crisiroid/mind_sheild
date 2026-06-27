package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type ConflictExerciseRepository struct {
	*BaseRepository
}

func NewConflictExerciseRepository(db *gorm.DB) *ConflictExerciseRepository {
	return &ConflictExerciseRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *ConflictExerciseRepository) Create(ctx context.Context, exercise *models.ConflictExercise) error {
	return r.BaseRepository.Create(ctx, exercise)
}

func (r *ConflictExerciseRepository) GetByID(ctx context.Context, id string) (*models.ConflictExercise, error) {
	var exercise models.ConflictExercise
	err := r.BaseRepository.GetByID(ctx, &exercise, id)
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *ConflictExerciseRepository) Update(ctx context.Context, exercise *models.ConflictExercise) error {
	return r.BaseRepository.Update(ctx, exercise)
}

func (r *ConflictExerciseRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.ConflictExercise{}, id)
}

func (r *ConflictExerciseRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.ConflictExercise, int64, error) {
	var exercises []models.ConflictExercise
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.ConflictExercise{})

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

func (r *ConflictExerciseRepository) GetByUserID(ctx context.Context, userID string) ([]models.ConflictExercise, error) {
	var exercises []models.ConflictExercise
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("scenario_id ASC").
		Find(&exercises).Error
	return exercises, err
}

func (r *ConflictExerciseRepository) GetExerciseStats(ctx context.Context, userID string) (*ConflictExerciseStats, error) {
	var stats ConflictExerciseStats

	query := r.DB.Model(&models.ConflictExercise{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.TotalScenarios)

	r.DB.Model(&models.ConflictExercise{}).
		Select("COALESCE(AVG(practice_count), 0)").
		Scan(&stats.AvgPracticeCount)

	r.DB.Model(&models.ConflictExercise{}).
		Where("performance_score IS NOT NULL").
		Select("COALESCE(AVG(performance_score), 0)").
		Scan(&stats.AvgPerformanceScore)

	return &stats, nil
}

func (r *ConflictExerciseRepository) GetPerformanceTrend(ctx context.Context, userID string) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			created_at as timestamp,
			scenario_id as count,
			performance_score as value
		FROM conflict_exercises
		WHERE user_id = ? AND performance_score IS NOT NULL
		ORDER BY scenario_id ASC
	`, userID).Scan(&trends)

	return trends, nil
}

func (r *ConflictExerciseRepository) GetScenarioCompletion(ctx context.Context, userID string) ([]ScenarioCompletionStats, error) {
	var completions []ScenarioCompletionStats

	query := r.DB.Raw(`
		SELECT
			scenario_id as label,
			practice_count as total,
			CASE WHEN performance_score IS NOT NULL THEN 1 ELSE 0 END as completed
		FROM conflict_exercises
		WHERE user_id = ? OR ? = ''
		ORDER BY scenario_id ASC
	`, userID, userID)

	query.Scan(&completions)
	return completions, nil
}

func (r *ConflictExerciseRepository) GetPracticeFrequency(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count,
			SUM(practice_count) as value
		FROM conflict_exercises
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

type ConflictExerciseStats struct {
	TotalScenarios      int64   `json:"total_scenarios"`
	AvgPracticeCount    float64 `json:"avg_practice_count"`
	AvgPerformanceScore float64 `json:"avg_performance_score"`
}

type ScenarioCompletionStats struct {
	Label     string `json:"label"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
}
