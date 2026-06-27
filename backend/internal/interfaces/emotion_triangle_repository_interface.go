package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type EmotionTriangleRepositoryInterface interface {
	Create(ctx context.Context, interaction *models.EmotionTriangleInteraction) error
	GetByID(ctx context.Context, id string) (*models.EmotionTriangleInteraction, error)
	Update(ctx context.Context, interaction *models.EmotionTriangleInteraction) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.EmotionTriangleInteraction, int64, error)
}
