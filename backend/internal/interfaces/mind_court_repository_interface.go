package interfaces

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type MindCourtRepositoryInterface interface {
	Create(ctx context.Context, evidence *models.MindCourtEvidence) error
	GetByID(ctx context.Context, id string) (*models.MindCourtEvidence, error)
	Update(ctx context.Context, evidence *models.MindCourtEvidence) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MindCourtEvidence, int64, error)
}
