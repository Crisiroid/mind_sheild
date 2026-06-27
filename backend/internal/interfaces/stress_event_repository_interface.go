package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type StressEventRepositoryInterface interface {
	Create(ctx context.Context, event *models.StressEvent) error
	GetByID(ctx context.Context, id string) (*models.StressEvent, error)
	Update(ctx context.Context, event *models.StressEvent) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.StressEvent, int64, error)
	GetStressStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.StressEventStats, error)
	GetIntensityTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetSituationTypeDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error)
}
