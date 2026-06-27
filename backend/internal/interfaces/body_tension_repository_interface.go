package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type BodyTensionRepositoryInterface interface {
	Create(ctx context.Context, tensionMap *models.BodyTensionMap) error
	GetByID(ctx context.Context, id string) (*models.BodyTensionMap, error)
	Update(ctx context.Context, tensionMap *models.BodyTensionMap) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.BodyTensionMap, int64, error)
	GetAverageIntensity(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*schemas.IntensityStatsResponse, error)
	GetIntensityTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetSeverityColorDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error)
}
