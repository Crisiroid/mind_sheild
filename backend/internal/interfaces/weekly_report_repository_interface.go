package interfaces

import (
	"context"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"

	"gorm.io/gorm"
)

type WeeklyReportRepositoryInterface interface {
	Create(ctx context.Context, report *models.WeeklyReport) error
	GetByID(ctx context.Context, id string) (*models.WeeklyReport, error)
	GetByUserAndWeek(ctx context.Context, userID string, weekNumber int) (*models.WeeklyReport, error)
	Update(ctx context.Context, report *models.WeeklyReport) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.WeeklyReport, int64, error)
	GetWeeklyStats(ctx context.Context, userID string) (*repository.WeeklyReportStats, error)
}
