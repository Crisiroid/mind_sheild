package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type WeeklyMediaContentRepositoryInterface interface {
	Create(ctx context.Context, media *models.WeeklyMediaContent) error
	GetByID(ctx context.Context, id string) (*models.WeeklyMediaContent, error)
	Update(ctx context.Context, media *models.WeeklyMediaContent) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.WeeklyMediaContent, int64, error)
	IncrementDownloadCount(ctx context.Context, id string) error
	GetByWeekNumber(ctx context.Context, weekNumber int, page, pageSize int) ([]models.WeeklyMediaContent, int64, error)
}
