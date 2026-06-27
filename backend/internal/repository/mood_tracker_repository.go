package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type MoodTrackerRepository struct {
	*BaseRepository
}

func NewMoodTrackerRepository(db *gorm.DB) *MoodTrackerRepository {
	return &MoodTrackerRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *MoodTrackerRepository) Create(ctx context.Context, tracker *models.MoodTracker) error {
	return r.BaseRepository.Create(ctx, tracker)
}

func (r *MoodTrackerRepository) GetByID(ctx context.Context, id string) (*models.MoodTracker, error) {
	var tracker models.MoodTracker
	err := r.BaseRepository.GetByID(ctx, &tracker, id)
	if err != nil {
		return nil, err
	}
	return &tracker, nil
}

func (r *MoodTrackerRepository) Update(ctx context.Context, tracker *models.MoodTracker) error {
	return r.BaseRepository.Update(ctx, tracker)
}

func (r *MoodTrackerRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.MoodTracker{}, id)
}

func (r *MoodTrackerRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MoodTracker, int64, error) {
	var trackers []models.MoodTracker
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.MoodTracker{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&trackers).Error; err != nil {
		return nil, 0, err
	}

	return trackers, total, nil
}

func (r *MoodTrackerRepository) GetByUserID(ctx context.Context, userID string) ([]models.MoodTracker, error) {
	var trackers []models.MoodTracker
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&trackers).Error
	return trackers, err
}

func (r *MoodTrackerRepository) GetMoodStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*MoodTrackerStats, error) {
	var stats MoodTrackerStats

	query := r.DB.Model(&models.MoodTracker{})

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

	r.DB.Model(&models.MoodTracker{}).
		Select("COALESCE(AVG(mood_before), 0)").
		Scan(&stats.AvgMoodBefore)

	r.DB.Model(&models.MoodTracker{}).
		Select("COALESCE(AVG(mood_after), 0)").
		Scan(&stats.AvgMoodAfter)

	if stats.AvgMoodBefore > 0 {
		stats.MoodImprovement = stats.AvgMoodAfter - stats.AvgMoodBefore
		stats.ImprovementPercentage = (stats.MoodImprovement / stats.AvgMoodBefore) * 100
	}

	return &stats, nil
}

func (r *MoodTrackerRepository) GetMoodTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]MoodTrendPoint, error) {
	var trends []MoodTrendPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			AVG(mood_before) as mood_before,
			AVG(mood_after) as mood_after
		FROM mood_tracker
		WHERE user_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *MoodTrackerRepository) GetActivityEffectiveness(ctx context.Context, userID string) ([]ActivityEffectiveness, error) {
	var rankings []ActivityEffectiveness

	query := r.DB.Raw(`
		SELECT
			activity_name as label,
			COUNT(*) as count,
			AVG(mood_after - mood_before) as avg_improvement
		FROM mood_tracker
		WHERE user_id = ? OR ? = ''
		GROUP BY activity_name
		ORDER BY avg_improvement DESC
	`, userID, userID)

	query.Scan(&rankings)
	return rankings, nil
}

func (r *MoodTrackerRepository) GetMoodDistribution(ctx context.Context, userID string, isBefore bool) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	moodField := "mood_after"
	if isBefore {
		moodField = "mood_before"
	}

	query := r.DB.Model(&models.MoodTracker{}).
		Select(moodField + " as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group(moodField).
		Order(moodField + " ASC").
		Scan(&distributions)

	return distributions, nil
}

type MoodTrackerStats struct {
	TotalEntries          int64   `json:"total_entries"`
	AvgMoodBefore         float64 `json:"avg_mood_before"`
	AvgMoodAfter          float64 `json:"avg_mood_after"`
	MoodImprovement       float64 `json:"mood_improvement"`
	ImprovementPercentage float64 `json:"improvement_percentage"`
}

type MoodTrendPoint struct {
	Timestamp  time.Time `json:"timestamp"`
	MoodBefore float64   `json:"mood_before"`
	MoodAfter  float64   `json:"mood_after"`
}

type ActivityEffectiveness struct {
	Label          string  `json:"label"`
	Count          int64   `json:"count"`
	AvgImprovement float64 `json:"avg_improvement"`
}
