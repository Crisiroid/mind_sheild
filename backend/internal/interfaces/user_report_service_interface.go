package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type UserReportServiceInterface interface {
	CreateUserReport(ctx context.Context, req *schemas.UserReportCreateRequest) (*schemas.UserReportResponse, error)
	GetUserReportById(ctx context.Context, id string) (*schemas.UserReportResponse, error)
	GetAllUserReports(ctx context.Context, req *schemas.UserReportListRequest) (*schemas.UserReportListResponse, error)
	DeleteUserReport(ctx context.Context, id string) error
}
