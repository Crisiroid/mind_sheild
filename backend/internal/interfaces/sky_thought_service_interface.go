package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type SkyThoughtServiceInterface interface {
	CreateSkyThought(ctx context.Context, req *schemas.SkyThoughtCreateRequest) (*schemas.SkyThoughtResponse, error)
	GetSkyThoughtById(ctx context.Context, id string) (*schemas.SkyThoughtResponse, error)
	GetAllSkyThoughts(ctx context.Context, req *schemas.SkyThoughtListRequest) (*schemas.SkyThoughtListResponse, error)
	UpdateSkyThought(ctx context.Context, id string, req *schemas.SkyThoughtUpdateRequest) (*schemas.SkyThoughtResponse, error)
	DeleteSkyThought(ctx context.Context, id string) error
}
