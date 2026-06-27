package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type AdminRoleRepositoryInterface interface {
	Create(ctx context.Context, role *models.AdminRole) error
	GetByID(ctx context.Context, id string) (*models.AdminRole, error)
	GetByRoleName(ctx context.Context, roleName string) (*models.AdminRole, error)
	Update(ctx context.Context, role *models.AdminRole) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.AdminRole, int64, error)
}
