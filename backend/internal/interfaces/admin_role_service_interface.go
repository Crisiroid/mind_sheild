package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type AdminRoleServiceInterface interface {
	CreateAdminRole(ctx context.Context, req *schemas.AdminRoleCreateRequest) (*schemas.AdminRoleResponse, error)
	GetAdminRoleById(ctx context.Context, roleId string) (*schemas.AdminRoleResponse, error)
	GetAllRoles(ctx context.Context, req *schemas.AdminRoleListRequest) (*schemas.AdminRoleListResponse, error)
	UpdateAdminRole(ctx context.Context, id string, req *schemas.AdminRoleUpdateRequest) (*schemas.AdminRoleResponse, error)
	DeleteAdminRole(ctx context.Context, id string) error
}
