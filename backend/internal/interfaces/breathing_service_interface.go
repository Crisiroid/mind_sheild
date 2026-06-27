package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"
)

type BreathingServiceInterface interface {
	CreateBreathingSession(ctx context.Context, req *schemas.BreathingSessionCreateRequest) (*schemas.BreathingSessionResponse, error)
	GetBreathingSessionById(ctx context.Context, id string) (*schemas.BreathingSessionResponse, error)
	GetAllBreathingSessions(ctx context.Context, req *schemas.BreathingSessionListRequest) (*schemas.BreathingSessionListResponse, error)
	UpdateBreathingSession(ctx context.Context, id string, req *schemas.BreathingSessionUpdateRequest) (*schemas.BreathingSessionResponse, error)
	DeleteBreathingSession(ctx context.Context, id string) error
	GetSessionStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.BreathingSessionStats, error)
	GetDurationTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetPatternUsage(ctx context.Context, userID string) ([]schemas.DistributionStats, error)
}
