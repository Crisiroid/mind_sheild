package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type SystemLogRepositoryInterface interface {
	Create(ctx context.Context, log *models.SystemLog) error
	GetByID(ctx context.Context, id string) (*models.SystemLog, error)
	Update(ctx context.Context, log *models.SystemLog) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.SystemLog, int64, error)
}
