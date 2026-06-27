package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type UserReportRepositoryInterface interface {
	Create(ctx context.Context, report *models.UserReport) error
	GetByID(ctx context.Context, id string) (*models.UserReport, error)
	Update(ctx context.Context, report *models.UserReport) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.UserReport, int64, error)
}
