package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type NegativeThoughtServiceInterface interface {
	CreateNegativeThought(ctx context.Context, req *schemas.NegativeThoughtCreateRequest) (*schemas.NegativeThoughtResponse, error)
	GetNegativeThoughtById(ctx context.Context, id string) (*schemas.NegativeThoughtResponse, error)
	GetAllNegativeThoughts(ctx context.Context, req *schemas.NegativeThoughtListRequest) (*schemas.NegativeThoughtListResponse, error)
	UpdateNegativeThought(ctx context.Context, id string, req *schemas.NegativeThoughtUpdateRequest) (*schemas.NegativeThoughtResponse, error)
	DeleteNegativeThought(ctx context.Context, id string) error
}
