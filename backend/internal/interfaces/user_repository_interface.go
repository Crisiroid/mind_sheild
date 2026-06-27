package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.User, int64, error)
	GetUserStats(ctx context.Context) (*schemas.UserStatsResponse, error)
	GetUserActivityTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetLoginAnalytics(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetAgreementStats(ctx context.Context) (int64, int64, float64, error)
	GetAppVersionDistribution(ctx context.Context) ([]schemas.DistributionStats, error)
	GetInactiveUsers(ctx context.Context, daysThreshold int) ([]models.User, error)
	GetUserEngagementStats(ctx context.Context, dateFrom, dateTo time.Time) (*schemas.EngagementStatsResponse, error)
	ExportUsers(ctx context.Context, dateFrom, dateTo *time.Time, userID string) ([]models.User, error)
}
