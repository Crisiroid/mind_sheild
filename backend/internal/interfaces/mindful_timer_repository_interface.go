package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type MindfulTimerRepositoryInterface interface {
	Create(ctx context.Context, timer *models.MindfulTimer) error
	GetByID(ctx context.Context, id string) (*models.MindfulTimer, error)
	Update(ctx context.Context, timer *models.MindfulTimer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MindfulTimer, int64, error)
}
