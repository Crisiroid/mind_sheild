package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"
)

type CognitiveGameServiceInterface interface {
	CreateCognitiveGame(ctx context.Context, req *schemas.CognitiveGameCreateRequest) (*schemas.CognitiveGameResponse, error)
	GetCognitiveGameById(ctx context.Context, id string) (*schemas.CognitiveGameResponse, error)
	GetAllCognitiveGames(ctx context.Context, req *schemas.CognitiveGameListRequest) (*schemas.CognitiveGameListResponse, error)
	UpdateCognitiveGame(ctx context.Context, id string, req *schemas.CognitiveGameUpdateRequest) (*schemas.CognitiveGameResponse, error)
	DeleteCognitiveGame(ctx context.Context, id string) error
	GetGameStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.CognitiveGameStats, error)
	GetScoreTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetTimeAnalysis(ctx context.Context, userID string) (*schemas.TimeAnalysisResponse, error)
}
