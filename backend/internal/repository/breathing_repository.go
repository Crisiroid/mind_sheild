package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type BreathingRepository struct {
	*BaseRepository
}

func NewBreathingRepository(db *gorm.DB) *BreathingRepository {
	return &BreathingRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *BreathingRepository) Create(ctx context.Context, session *models.BreathingSession) error {
	return r.BaseRepository.Create(ctx, session)
}

func (r *BreathingRepository) GetByID(ctx context.Context, id string) (*models.BreathingSession, error) {
	var session models.BreathingSession
	err := r.BaseRepository.GetByID(ctx, &session, id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *BreathingRepository) Update(ctx context.Context, session *models.BreathingSession) error {
	return r.BaseRepository.Update(ctx, session)
}

func (r *BreathingRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.BreathingSession{}, id)
}

func (r *BreathingRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.BreathingSession, int64, error) {
	var sessions []models.BreathingSession
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.BreathingSession{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

func (r *BreathingRepository) GetByUserID(ctx context.Context, userID string) ([]models.BreathingSession, error) {
	var sessions []models.BreathingSession
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

func (r *BreathingRepository) GetSessionStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*BreathingSessionStats, error) {
	var stats BreathingSessionStats

	query := r.DB.Model(&models.BreathingSession{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if dateFrom != nil {
		query = query.Where("session_start >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("session_start <= ?", *dateTo)
	}

	query.Count(&stats.TotalSessions)

	r.DB.Model(&models.BreathingSession{}).
		Where("is_completed = ?", true).
		Count(&stats.CompletedSessions)

	stats.IncompleteSessions = stats.TotalSessions - stats.CompletedSessions
	if stats.TotalSessions > 0 {
		stats.CompletionRate = float64(stats.CompletedSessions) / float64(stats.TotalSessions) * 100
	}

	r.DB.Model(&models.BreathingSession{}).
		Select("COALESCE(AVG(duration_seconds), 0)").
		Scan(&stats.AvgDuration)

	return &stats, nil
}

func (r *BreathingRepository) GetDurationTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(session_start) as timestamp,
			COUNT(*) as count,
			AVG(duration_seconds) as value
		FROM breathing_sessions
		WHERE user_id = ? AND session_start >= ? AND session_start <= ?
		GROUP BY DATE(session_start)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *BreathingRepository) GetPatternUsage(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.BreathingSession{}).
		Select("breathing_pattern as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("breathing_pattern").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *BreathingRepository) GetCompletionRate(ctx context.Context, userID string) ([]PatternCompletionRate, error) {
	var rates []PatternCompletionRate

	query := r.DB.Raw(`
		SELECT
			breathing_pattern as label,
			COUNT(*) as total,
			SUM(CASE WHEN is_completed = true THEN 1 ELSE 0 END) as completed
		FROM breathing_sessions
		WHERE user_id = ? OR ? = ''
		GROUP BY breathing_pattern
		ORDER BY total DESC
	`, userID, userID)

	query.Scan(&rates)

	for i := range rates {
		if rates[i].Total > 0 {
			rates[i].Rate = float64(rates[i].Completed) / float64(rates[i].Total) * 100
		}
	}

	return rates, nil
}

func (r *BreathingRepository) GetCalendarTickStats(ctx context.Context, userID string, dateFrom, dateTo time.Time) (int64, int64, float64, error) {
	var total, ticked int64

	query := r.DB.Model(&models.BreathingSession{}).Where("user_id = ?", userID)
	if !dateFrom.IsZero() {
		query = query.Where("session_start >= ?", dateFrom)
	}
	if !dateTo.IsZero() {
		query = query.Where("session_start <= ?", dateTo)
	}

	query.Count(&total)
	r.DB.Model(&models.BreathingSession{}).
		Where("user_id = ? AND calendar_ticked = ?", userID, true).
		Count(&ticked)

	var rate float64
	if total > 0 {
		rate = float64(ticked) / float64(total) * 100
	}

	return total, ticked, rate, nil
}

type BreathingSessionStats struct {
	TotalSessions      int64   `json:"total_sessions"`
	CompletedSessions  int64   `json:"completed_sessions"`
	IncompleteSessions int64   `json:"incomplete_sessions"`
	CompletionRate     float64 `json:"completion_rate"`
	AvgDuration        float64 `json:"avg_duration"`
}

type PatternCompletionRate struct {
	Label     string  `json:"label"`
	Total     int64   `json:"total"`
	Completed int64   `json:"completed"`
	Rate      float64 `json:"rate"`
}
