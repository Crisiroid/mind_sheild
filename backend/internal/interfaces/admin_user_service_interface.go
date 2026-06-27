package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type AdminUserServiceInterface interface {
	CreateAdminUser(ctx context.Context, req *schemas.AdminCreateRequest) (*schemas.AdminUserResponse, error)
	GetAdminUserById(ctx context.Context, adminUserId string) (*schemas.AdminUserResponse, error)
	GetAdminUserByUsername(ctx context.Context, adminUserName string) (*schemas.AdminUserResponse, error)
	GetAdminUserByEmail(ctx context.Context, adminUserEmail string) (*schemas.AdminUserResponse, error)
	GetAllAdmins(ctx context.Context, req *schemas.AdminListRequest) (*schemas.AdminListResponse, error)
	UpdateAdminUser(ctx context.Context, id string, req *schemas.AdminUpdateRequest) (*schemas.AdminUserResponse, error)
	DeleteAdminUser(ctx context.Context, id string) error
	DeactivateAdminUser(ctx context.Context, id string) (*schemas.AdminUserResponse, error)
	GetAdminProfile(ctx context.Context, adminID string) (*schemas.AdminUserResponse, error)
	UpdateAdminProfile(ctx context.Context, adminID string, req *schemas.AdminUpdateProfileRequest) (*schemas.AdminUserResponse, error)
}
