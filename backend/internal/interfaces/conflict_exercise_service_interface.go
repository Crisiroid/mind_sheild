package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type ConflictExerciseServiceInterface interface {
	CreateConflictExercise(ctx context.Context, req *schemas.ConflictExerciseCreateRequest) (*schemas.ConflictExerciseResponse, error)
	GetConflictExerciseById(ctx context.Context, id string) (*schemas.ConflictExerciseResponse, error)
	GetAllConflictExercises(ctx context.Context, req *schemas.ConflictExerciseListRequest) (*schemas.ConflictExerciseListResponse, error)
	UpdateConflictExercise(ctx context.Context, id string, req *schemas.ConflictExerciseUpdateRequest) (*schemas.ConflictExerciseResponse, error)
	DeleteConflictExercise(ctx context.Context, id string) error
}
