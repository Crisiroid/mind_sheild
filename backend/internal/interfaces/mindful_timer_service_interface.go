package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type MindfulTimerServiceInterface interface {
	CreateMindfulTimer(ctx context.Context, req *schemas.MindfulTimerCreateRequest) (*schemas.MindfulTimerResponse, error)
	GetMindfulTimerById(ctx context.Context, id string) (*schemas.MindfulTimerResponse, error)
	GetAllMindfulTimers(ctx context.Context, req *schemas.MindfulTimerListRequest) (*schemas.MindfulTimerListResponse, error)
	UpdateMindfulTimer(ctx context.Context, id string, req *schemas.MindfulTimerUpdateRequest) (*schemas.MindfulTimerResponse, error)
	DeleteMindfulTimer(ctx context.Context, id string) error
}
