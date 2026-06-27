package interfaces

import (
	"context"

	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"
)

type WeeklyReportServiceInterface interface {
	CreateWeeklyReport(ctx context.Context, req *schemas.WeeklyReportCreateRequest) (*schemas.WeeklyReportResponse, error)
	GetWeeklyReportById(ctx context.Context, id string) (*schemas.WeeklyReportResponse, error)
	GetAllWeeklyReports(ctx context.Context, req *schemas.WeeklyReportListRequest) (*schemas.WeeklyReportListResponse, error)
	UpdateWeeklyReport(ctx context.Context, id string, req *schemas.WeeklyReportUpdateRequest) (*schemas.WeeklyReportResponse, error)
	DeleteWeeklyReport(ctx context.Context, id string) error
	GetWeeklyStats(ctx context.Context, userID string) (*repository.WeeklyReportStats, error)
}
