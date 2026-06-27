package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type CognitiveGameRepository struct {
	*BaseRepository
}

func NewCognitiveGameRepository(db *gorm.DB) *CognitiveGameRepository {
	return &CognitiveGameRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *CognitiveGameRepository) Create(ctx context.Context, game *models.CognitiveErrorGame) error {
	return r.BaseRepository.Create(ctx, game)
}

func (r *CognitiveGameRepository) GetByID(ctx context.Context, id string) (*models.CognitiveErrorGame, error) {
	var game models.CognitiveErrorGame
	err := r.BaseRepository.GetByID(ctx, &game, id)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *CognitiveGameRepository) Update(ctx context.Context, game *models.CognitiveErrorGame) error {
	return r.BaseRepository.Update(ctx, game)
}

func (r *CognitiveGameRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.CognitiveErrorGame{}, id)
}

func (r *CognitiveGameRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.CognitiveErrorGame, int64, error) {
	var games []models.CognitiveErrorGame
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.CognitiveErrorGame{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&games).Error; err != nil {
		return nil, 0, err
	}

	return games, total, nil
}

func (r *CognitiveGameRepository) GetByUserID(ctx context.Context, userID string) ([]models.CognitiveErrorGame, error) {
	var games []models.CognitiveErrorGame
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&games).Error
	return games, err
}

func (r *CognitiveGameRepository) GetGameStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*CognitiveGameStats, error) {
	var stats CognitiveGameStats

	query := r.DB.Model(&models.CognitiveErrorGame{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}

	query.Count(&stats.TotalGames)

	r.DB.Model(&models.CognitiveErrorGame{}).
		Where("is_correct = ?", true).
		Count(&stats.CorrectAnswers)

	r.DB.Model(&models.CognitiveErrorGame{}).
		Select("COALESCE(AVG(score), 0)").
		Scan(&stats.AvgScore)

	r.DB.Model(&models.CognitiveErrorGame{}).
		Select("COALESCE(AVG(time_taken_seconds), 0)").
		Scan(&stats.AvgTimeTaken)

	if stats.TotalGames > 0 {
		stats.AccuracyRate = float64(stats.CorrectAnswers) / float64(stats.TotalGames) * 100
	}

	return &stats, nil
}

func (r *CognitiveGameRepository) GetScoreTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count,
			AVG(score) as value
		FROM cognitive_error_games
		WHERE user_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *CognitiveGameRepository) GetScenarioTypePerformance(ctx context.Context, userID string) ([]ScenarioPerformance, error) {
	var performances []ScenarioPerformance

	query := r.DB.Raw(`
		SELECT
			scenario_type as label,
			COUNT(*) as total,
			SUM(CASE WHEN is_correct = true THEN 1 ELSE 0 END) as correct,
			AVG(score) as avg_score
		FROM cognitive_error_games
		WHERE user_id = ? OR ? = ''
		GROUP BY scenario_type
		ORDER BY total DESC
	`, userID, userID)

	query.Scan(&performances)

	for i := range performances {
		if performances[i].Total > 0 {
			performances[i].Accuracy = float64(performances[i].Correct) / float64(performances[i].Total) * 100
		}
	}

	return performances, nil
}

func (r *CognitiveGameRepository) GetAccuracyByScenario(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Raw(`
		SELECT
			scenario_type as label,
			SUM(CASE WHEN is_correct = true THEN 1 ELSE 0 END) as count
		FROM cognitive_error_games
		WHERE user_id = ? OR ? = ''
		GROUP BY scenario_type
		ORDER BY count DESC
	`, userID, userID)

	query.Scan(&distributions)

	return distributions, nil
}

func (r *CognitiveGameRepository) GetTimeAnalysis(ctx context.Context, userID string) (*schemas.TimeAnalysisResponse, error) {
	var stats schemas.TimeAnalysisResponse

	query := r.DB.Model(&models.CognitiveErrorGame{}).
		Where("time_taken_seconds IS NOT NULL")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.EntriesCount)
	r.DB.Model(&models.CognitiveErrorGame{}).
		Select("COALESCE(AVG(time_taken_seconds), 0)").
		Scan(&stats.AvgDuration)
	r.DB.Model(&models.CognitiveErrorGame{}).
		Select("COALESCE(MIN(time_taken_seconds), 0)").
		Scan(&stats.MinDuration)
	r.DB.Model(&models.CognitiveErrorGame{}).
		Select("COALESCE(MAX(time_taken_seconds), 0)").
		Scan(&stats.MaxDuration)

	return &stats, nil
}

type CognitiveGameStats struct {
	TotalGames     int64   `json:"total_games"`
	CorrectAnswers int64   `json:"correct_answers"`
	AccuracyRate   float64 `json:"accuracy_rate"`
	AvgScore       float64 `json:"avg_score"`
	AvgTimeTaken   float64 `json:"avg_time_taken"`
}

type ScenarioPerformance struct {
	Label    string  `json:"label"`
	Total    int64   `json:"total"`
	Correct  int64   `json:"correct"`
	AvgScore float64 `json:"avg_score"`
	Accuracy float64 `json:"accuracy"`
}
