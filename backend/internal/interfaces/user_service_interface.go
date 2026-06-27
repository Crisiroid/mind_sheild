package interfaces

import (
	"context"
	"time"

	"psychology-backend/pkg/schemas"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, req *schemas.UserCreateRequest) (*schemas.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*schemas.UserResponse, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*schemas.UserResponse, error)
	ListUsers(ctx context.Context, req *schemas.UserListRequest) (*schemas.UserListResponse, error)
	UpdateUser(ctx context.Context, id string, req *schemas.UserUpdateRequest) (*schemas.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
	AcceptAgreement(ctx context.Context, userID string) (*schemas.UserResponse, error)
	UpdateLoginInfo(ctx context.Context, userID string, androidVersion, appVersion string) (*schemas.UserResponse, error)
	GetUserStats(ctx context.Context) (*schemas.UserStatsResponse, error)
	GetUserActivityTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetLoginAnalytics(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error)
	GetAgreementStats(ctx context.Context) (int64, int64, float64, error)
	GetAppVersionDistribution(ctx context.Context) ([]schemas.DistributionStats, error)
	GetInactiveUsers(ctx context.Context, daysThreshold int) ([]schemas.UserResponse, error)
	GetUserEngagementStats(ctx context.Context, dateFrom, dateTo time.Time) (*schemas.EngagementStatsResponse, error)
	ExportUsers(ctx context.Context, dateFrom, dateTo *time.Time, userID string) ([]schemas.UserResponse, error)
	GetUserProfile(ctx context.Context, userID string) (*schemas.UserResponse, error)
	UpdateUserProfile(ctx context.Context, userID string, req *schemas.UserUpdateProfileRequest) (*schemas.UserResponse, error)
	SyncUserData(ctx context.Context, userID string, syncData map[string]interface{}) error
}
