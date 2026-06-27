package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type AcceptanceRepositoryInterface interface {
	Create(ctx context.Context, exercise *models.AcceptanceExercise) error
	GetByID(ctx context.Context, id string) (*models.AcceptanceExercise, error)
	Update(ctx context.Context, exercise *models.AcceptanceExercise) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.AcceptanceExercise, int64, error)
}
