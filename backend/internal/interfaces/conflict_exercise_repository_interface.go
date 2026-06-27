package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type ConflictExerciseRepositoryInterface interface {
	Create(ctx context.Context, exercise *models.ConflictExercise) error
	GetByID(ctx context.Context, id string) (*models.ConflictExercise, error)
	Update(ctx context.Context, exercise *models.ConflictExercise) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.ConflictExercise, int64, error)
}
