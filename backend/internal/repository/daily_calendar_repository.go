package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type DailyCalendarRepository struct {
	*BaseRepository
}

func NewDailyCalendarRepository(db *gorm.DB) *DailyCalendarRepository {
	return &DailyCalendarRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *DailyCalendarRepository) Create(ctx context.Context, calendar *models.DailyCalendar) error {
	return r.BaseRepository.Create(ctx, calendar)
}

func (r *DailyCalendarRepository) GetByID(ctx context.Context, id string) (*models.DailyCalendar, error) {
	var calendar models.DailyCalendar
	err := r.BaseRepository.GetByID(ctx, &calendar, id)
	if err != nil {
		return nil, err
	}
	return &calendar, nil
}

func (r *DailyCalendarRepository) GetByUserAndDay(ctx context.Context, userID string, dayNumber int) (*models.DailyCalendar, error) {
	var calendar models.DailyCalendar
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND day_number = ?", userID, dayNumber).
		First(&calendar).Error
	if err != nil {
		return nil, err
	}
	return &calendar, nil
}

func (r *DailyCalendarRepository) Update(ctx context.Context, calendar *models.DailyCalendar) error {
	return r.BaseRepository.Update(ctx, calendar)
}

func (r *DailyCalendarRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.DailyCalendar{}, id)
}

func (r *DailyCalendarRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.DailyCalendar, int64, error) {
	var calendars []models.DailyCalendar
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.DailyCalendar{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("day_number ASC").Find(&calendars).Error; err != nil {
		return nil, 0, err
	}

	return calendars, total, nil
}

func (r *DailyCalendarRepository) GetByUserID(ctx context.Context, userID string) ([]models.DailyCalendar, error) {
	var calendars []models.DailyCalendar
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("day_number ASC").
		Find(&calendars).Error
	return calendars, err
}

func (r *DailyCalendarRepository) GetCompletionStats(ctx context.Context, userID string) (*schemas.CompletionStatsResponse, error) {
	var stats schemas.CompletionStatsResponse

	query := r.DB.Model(&models.DailyCalendar{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.Total)
	r.DB.Model(&models.DailyCalendar{}).
		Where("is_completed = ?", true).
		Count(&stats.Completed)

	stats.Incomplete = stats.Total - stats.Completed
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	return &stats, nil
}

func (r *DailyCalendarRepository) GetDayRangeProgress(ctx context.Context, userID string, fromDay, toDay int) (*schemas.CompletionStatsResponse, error) {
	var stats schemas.CompletionStatsResponse

	query := r.DB.Model(&models.DailyCalendar{}).
		Where("day_number >= ? AND day_number <= ?", fromDay, toDay)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.Total)
	r.DB.Model(&models.DailyCalendar{}).
		Where("day_number >= ? AND day_number <= ? AND is_completed = ?", fromDay, toDay, true).
		Count(&stats.Completed)

	stats.Incomplete = stats.Total - stats.Completed
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	return &stats, nil
}

func (r *DailyCalendarRepository) GetActivityDistribution(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	r.DB.Raw(`
		SELECT
			DATE(calendar_date) as label,
			COUNT(*) as count
		FROM daily_calendar
		WHERE is_completed = true
			AND calendar_date >= ? AND calendar_date <= ?
		GROUP BY DATE(calendar_date)
		ORDER BY label ASC
	`, dateFrom, dateTo).Scan(&distributions)

	return distributions, nil
}

func (r *DailyCalendarRepository) GetStreakAnalysis(ctx context.Context, userID string) (*StreakAnalysis, error) {
	var analysis StreakAnalysis

	var calendars []models.DailyCalendar
	r.DB.WithContext(ctx).
		Where("user_id = ? AND is_completed = ?", userID, true).
		Order("day_number ASC").
		Find(&calendars)

	if len(calendars) == 0 {
		return &analysis, nil
	}

	currentStreak := 0
	maxStreak := 0
	tempStreak := 0

	for i, cal := range calendars {
		if i == 0 {
			tempStreak = 1
		} else {
			if cal.DayNumber == calendars[i-1].DayNumber+1 {
				tempStreak++
			} else {
				if tempStreak > maxStreak {
					maxStreak = tempStreak
				}
				tempStreak = 1
			}
		}
	}

	if tempStreak > maxStreak {
		maxStreak = tempStreak
	}

	analysis.CurrentStreak = currentStreak
	analysis.LongestStreak = maxStreak
	analysis.TotalCompleted = len(calendars)

	return &analysis, nil
}

func (r *DailyCalendarRepository) GetCompletionTrend(ctx context.Context, userID string) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	query := r.DB.Raw(`
		SELECT
			DATE(calendar_date) as timestamp,
			COUNT(*) as count
		FROM daily_calendar
		WHERE user_id = ? AND is_completed = true
		GROUP BY DATE(calendar_date)
		ORDER BY timestamp ASC
	`, userID)

	query.Scan(&trends)
	return trends, nil
}

type StreakAnalysis struct {
	CurrentStreak  int `json:"current_streak"`
	LongestStreak  int `json:"longest_streak"`
	TotalCompleted int `json:"total_completed"`
}
