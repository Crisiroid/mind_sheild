package interfaces

import (
	"context"
	"time"

	"psychology-backend/pkg/schemas"
)

type BodyTensionServiceInterface interface {
	CreateBodyTension(ctx context.Context, req *schemas.BodyTensionCreateRequest) (*schemas.BodyTensionResponse, error)
	GetBodyTensionById(ctx context.Context, id string) (*schemas.BodyTensionResponse, error)
	GetAllBodyTensions(ctx context.Context, req *schemas.BodyTensionListRequest) (*schemas.BodyTensionListResponse, error)
	UpdateBodyTension(ctx context.Context, id string, req *schemas.BodyTensionUpdateRequest) (*schemas.BodyTensionResponse, error)
	DeleteBodyTension(ctx context.Context, id string) error
	GetAverageIntensity(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*schemas.IntensityStatsResponse, error)
	GetIntensityTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetSeverityColorDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error)
}
