package interfaces

import (
	"context"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type DailyCalendarRepositoryInterface interface {
	Create(ctx context.Context, calendar *models.DailyCalendar) error
	GetByID(ctx context.Context, id string) (*models.DailyCalendar, error)
	GetByUserAndDay(ctx context.Context, userID string, dayNumber int) (*models.DailyCalendar, error)
	Update(ctx context.Context, calendar *models.DailyCalendar) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.DailyCalendar, int64, error)
	GetCompletionStats(ctx context.Context, userID string) (*schemas.CompletionStatsResponse, error)
	GetDayRangeProgress(ctx context.Context, userID string, fromDay, toDay int) (*schemas.CompletionStatsResponse, error)
	GetStreakAnalysis(ctx context.Context, userID string) (*repository.StreakAnalysis, error)
}
