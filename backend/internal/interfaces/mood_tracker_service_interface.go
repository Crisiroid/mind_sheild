package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"
)

type MoodTrackerServiceInterface interface {
	CreateMoodTracker(ctx context.Context, req *schemas.MoodTrackerCreateRequest) (*schemas.MoodTrackerResponse, error)
	GetMoodTrackerById(ctx context.Context, id string) (*schemas.MoodTrackerResponse, error)
	GetAllMoodTrackers(ctx context.Context, req *schemas.MoodTrackerListRequest) (*schemas.MoodTrackerListResponse, error)
	UpdateMoodTracker(ctx context.Context, id string, req *schemas.MoodTrackerUpdateRequest) (*schemas.MoodTrackerResponse, error)
	DeleteMoodTracker(ctx context.Context, id string) error
	GetMoodStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.MoodTrackerStats, error)
	GetMoodTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]repository.MoodTrendPoint, error)
	GetActivityEffectiveness(ctx context.Context, userID string) ([]repository.ActivityEffectiveness, error)
}
