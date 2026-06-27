package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type MentalMustRepositoryInterface interface {
	Create(ctx context.Context, mentalMust *models.MentalMust) error
	GetByID(ctx context.Context, id string) (*models.MentalMust, error)
	Update(ctx context.Context, mentalMust *models.MentalMust) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MentalMust, int64, error)
}
