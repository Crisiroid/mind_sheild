package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type MindCourtServiceInterface interface {
	CreateMindCourt(ctx context.Context, req *schemas.MindCourtCreateRequest) (*schemas.MindCourtResponse, error)
	GetMindCourtById(ctx context.Context, id string) (*schemas.MindCourtResponse, error)
	GetAllMindCourts(ctx context.Context, req *schemas.MindCourtListRequest) (*schemas.MindCourtListResponse, error)
	UpdateMindCourt(ctx context.Context, id string, req *schemas.MindCourtUpdateRequest) (*schemas.MindCourtResponse, error)
	DeleteMindCourt(ctx context.Context, id string) error
}
