package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type RoleValueRepositoryInterface interface {
	Create(ctx context.Context, roleValue *models.RoleAndValue) error
	GetByID(ctx context.Context, id string) (*models.RoleAndValue, error)
	Update(ctx context.Context, roleValue *models.RoleAndValue) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.RoleAndValue, int64, error)
}
