package interfaces

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"

	"gorm.io/gorm"
)

type AdminUserRepositoryInterface interface {
	Create(ctx context.Context, adminUser *models.AdminUser) error
	GetByID(ctx context.Context, id string) (*models.AdminUser, error)
	GetByUsername(ctx context.Context, username string) (*models.AdminUser, error)
	GetByEmail(ctx context.Context, email string) (*models.AdminUser, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*models.AdminUser, error)
	Update(ctx context.Context, adminUser *models.AdminUser) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.AdminUser, int64, error)
	GetAdminStats(ctx context.Context) (int64, int64, error)
	GetAdminActivityLog(ctx context.Context, dateFrom, dateTo time.Time) ([]repository.AdminActivityLog, error)
	GetAdminsByRole(ctx context.Context) ([]repository.RoleDistribution, error)
	GetInactiveAdmins(ctx context.Context, daysThreshold int) ([]models.AdminUser, error)
	GetAdminsByActiveStatus(ctx context.Context, isActive bool) ([]models.AdminUser, error)
}
