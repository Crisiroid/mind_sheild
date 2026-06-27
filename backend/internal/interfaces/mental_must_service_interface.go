package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type MentalMustServiceInterface interface {
	CreateMentalMust(ctx context.Context, req *schemas.MentalMustCreateRequest) (*schemas.MentalMustResponse, error)
	GetMentalMustById(ctx context.Context, id string) (*schemas.MentalMustResponse, error)
	GetAllMentalMusts(ctx context.Context, req *schemas.MentalMustListRequest) (*schemas.MentalMustListResponse, error)
	UpdateMentalMust(ctx context.Context, id string, req *schemas.MentalMustUpdateRequest) (*schemas.MentalMustResponse, error)
	DeleteMentalMust(ctx context.Context, id string) error
}
