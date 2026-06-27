package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type BreathingRepositoryInterface interface {
	Create(ctx context.Context, session *models.BreathingSession) error
	GetByID(ctx context.Context, id string) (*models.BreathingSession, error)
	Update(ctx context.Context, session *models.BreathingSession) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.BreathingSession, int64, error)
	GetSessionStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.BreathingSessionStats, error)
	GetDurationTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetPatternUsage(ctx context.Context, userID string) ([]schemas.DistributionStats, error)
}
