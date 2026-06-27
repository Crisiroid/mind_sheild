package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type SystemLogServiceInterface interface {
	GetSystemLogById(ctx context.Context, id string) (*schemas.SystemLogResponse, error)
	GetAllSystemLogs(ctx context.Context, req *schemas.SystemLogListRequest) (*schemas.SystemLogListResponse, error)
	DeleteSystemLog(ctx context.Context, id string) error
}
