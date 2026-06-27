package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type CognitiveGameRepositoryInterface interface {
	Create(ctx context.Context, game *models.CognitiveErrorGame) error
	GetByID(ctx context.Context, id string) (*models.CognitiveErrorGame, error)
	Update(ctx context.Context, game *models.CognitiveErrorGame) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.CognitiveErrorGame, int64, error)
	GetGameStats(ctx context.Context, userID string, dateFrom, dateTo *time.Time) (*repository.CognitiveGameStats, error)
	GetScoreTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetTimeAnalysis(ctx context.Context, userID string) (*schemas.TimeAnalysisResponse, error)
}
