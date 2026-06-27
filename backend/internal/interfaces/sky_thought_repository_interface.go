package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type SkyThoughtRepositoryInterface interface {
	Create(ctx context.Context, skyThought *models.SkyThought) error
	GetByID(ctx context.Context, id string) (*models.SkyThought, error)
	Update(ctx context.Context, skyThought *models.SkyThought) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.SkyThought, int64, error)
}
