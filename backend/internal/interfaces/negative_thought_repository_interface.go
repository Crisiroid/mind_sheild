package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type NegativeThoughtRepositoryInterface interface {
	Create(ctx context.Context, thought *models.NegativeThought) error
	GetByID(ctx context.Context, id string) (*models.NegativeThought, error)
	Update(ctx context.Context, thought *models.NegativeThought) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.NegativeThought, int64, error)
}
