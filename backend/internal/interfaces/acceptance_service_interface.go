package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type AcceptanceServiceInterface interface {
	CreateAcceptanceExercise(ctx context.Context, req *schemas.AcceptanceExerciseCreateRequest) (*schemas.AcceptanceExerciseResponse, error)
	GetAcceptanceExerciseById(ctx context.Context, id string) (*schemas.AcceptanceExerciseResponse, error)
	GetAllAcceptanceExercises(ctx context.Context, req *schemas.AcceptanceExerciseListRequest) (*schemas.AcceptanceExerciseListResponse, error)
	UpdateAcceptanceExercise(ctx context.Context, id string, req *schemas.AcceptanceExerciseUpdateRequest) (*schemas.AcceptanceExerciseResponse, error)
	DeleteAcceptanceExercise(ctx context.Context, id string) error
}
