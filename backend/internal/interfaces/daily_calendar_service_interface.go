package interfaces

import (
	"context"

	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"
)

type DailyCalendarServiceInterface interface {
	CreateCalendarEntry(ctx context.Context, req *schemas.CalendarCreateRequest) (*schemas.CalendarResponse, error)
	GetCalendarEntryById(ctx context.Context, id string) (*schemas.CalendarResponse, error)
	GetAllCalendarEntries(ctx context.Context, req *schemas.CalendarListRequest) (*schemas.CalendarListResponse, error)
	UpdateCalendarEntry(ctx context.Context, id string, req *schemas.CalendarUpdateRequest) (*schemas.CalendarResponse, error)
	DeleteCalendarEntry(ctx context.Context, id string) error
	GetCompletionStats(ctx context.Context, userID string) (*schemas.CompletionStatsResponse, error)
	GetDayRangeProgress(ctx context.Context, userID string, fromDay, toDay int) (*schemas.CompletionStatsResponse, error)
	GetStreakAnalysis(ctx context.Context, userID string) (*repository.StreakAnalysis, error)
}
