package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type RoleValueServiceInterface interface {
	CreateRoleValue(ctx context.Context, req *schemas.RoleValueCreateRequest) (*schemas.RoleValueResponse, error)
	GetRoleValueById(ctx context.Context, id string) (*schemas.RoleValueResponse, error)
	GetAllRoleValues(ctx context.Context, req *schemas.RoleValueListRequest) (*schemas.RoleValueListResponse, error)
	UpdateRoleValue(ctx context.Context, id string, req *schemas.RoleValueUpdateRequest) (*schemas.RoleValueResponse, error)
	DeleteRoleValue(ctx context.Context, id string) error
}
