package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"

	"gorm.io/gorm"
)

type MoodTrackerRepositoryInterface interface {
	Create(ctx context.Context, tracker *models.MoodTracker) error
	GetByID(ctx context.Context, id string) (*models.MoodTracker, error)
	Update(ctx context.Context, tracker *models.MoodTracker) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MoodTracker, int64, error)
	GetMoodStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.MoodTrackerStats, error)
	GetMoodTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]repository.MoodTrendPoint, error)
	GetActivityEffectiveness(ctx context.Context, userID string) ([]repository.ActivityEffectiveness, error)
}
