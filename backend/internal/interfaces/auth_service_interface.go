package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type AuthServiceInterface interface {
	UserRegister(ctx context.Context, req *schemas.UserRegisterRequest) (*schemas.UserLoginResponse, error)
	UserLogin(ctx context.Context, req *schemas.UserLoginRequest) (*schemas.UserLoginResponse, error)
	UserRefreshToken(ctx context.Context, refreshToken string) (*schemas.UserRefreshTokenResponse, error)
	UserLogout(ctx context.Context, userID string) error
	UserChangePassword(ctx context.Context, userID string, req *schemas.UserChangePasswordRequest) error
	AdminLogin(ctx context.Context, req *schemas.AdminLoginRequest) (*schemas.AdminAuthLoginResponse, error)
	AdminRefreshToken(ctx context.Context, refreshToken string) (*schemas.AdminRefreshTokenResponse, error)
	AdminLogout(ctx context.Context, adminID string) error
	AdminChangePassword(ctx context.Context, adminID string, req *schemas.AdminChangePasswordRequest) error
}
