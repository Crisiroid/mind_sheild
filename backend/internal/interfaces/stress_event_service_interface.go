package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"
)

type StressEventServiceInterface interface {
	CreateStressEvent(ctx context.Context, req *schemas.StressEventCreateRequest) (*schemas.StressEventResponse, error)
	GetStressEventById(ctx context.Context, id string) (*schemas.StressEventResponse, error)
	GetAllStressEvents(ctx context.Context, req *schemas.StressEventListRequest) (*schemas.StressEventListResponse, error)
	UpdateStressEvent(ctx context.Context, id string, req *schemas.StressEventUpdateRequest) (*schemas.StressEventResponse, error)
	DeleteStressEvent(ctx context.Context, id string) error
	GetStressStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.StressEventStats, error)
	GetIntensityTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetSituationTypeDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error)
}
