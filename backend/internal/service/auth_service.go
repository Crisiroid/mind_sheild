package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo        interfaces.UserRepositoryInterface
	adminRepo       interfaces.AdminUserRepositoryInterface
	jwtService      *JWTService
	passwordService *PasswordService
}

func NewAuthService(
	userRepo interfaces.UserRepositoryInterface,
	adminRepo interfaces.AdminUserRepositoryInterface,
	jwtService *JWTService,
	passwordService *PasswordService,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		adminRepo:       adminRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

func (s *AuthService) UserRegister(ctx context.Context, req *schemas.UserRegisterRequest) (*schemas.UserLoginResponse, error) {
	if err := s.validatePhoneNumber(req.PhoneNumber); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existingUser, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingUser != nil {
		return nil, errors.New("کاربر با این شماره تلفن قبلاً ثبت نام کرده است")
	}

	passwordHash, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	user := &models.User{
		PhoneNumber:      req.PhoneNumber,
		PasswordHash:     passwordHash,
		RegistrationDate: now,
		LastLogin:        &now,
		LoginCount:       1,
		AndroidVersion:   req.AndroidVersion,
		AppVersion:       req.AppVersion,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.generateUserTokens(ctx, user)
}

func (s *AuthService) UserLogin(ctx context.Context, req *schemas.UserLoginRequest) (*schemas.UserLoginResponse, error) {
	user, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("شماره تلفن یا رمز عبور اشتباه است")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.passwordService.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.New("شماره تلفن یا رمز عبور اشتباه است")
	}

	now := time.Now()
	user.LastLogin = &now
	user.LoginCount++
	if err := s.userRepo.Update(ctx, user); err != nil {
		fmt.Printf("Warning: failed to update login info for user %s: %v\n", user.ID, err)
	}

	return s.generateUserTokens(ctx, user)
}

func (s *AuthService) UserRefreshToken(ctx context.Context, refreshToken string) (*schemas.UserRefreshTokenResponse, error) {
	user, err := s.userRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token expired or invalid")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID.String(), "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, expiry, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	user.RefreshToken = newRefreshToken
	user.RefreshTokenExpiry = &expiry
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &schemas.UserRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.jwtService.GetAccessExpiration().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthService) UserLogout(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.RefreshToken = ""
	user.RefreshTokenExpiry = nil

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

func (s *AuthService) UserChangePassword(ctx context.Context, userID string, req *schemas.UserChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.passwordService.CheckPassword(req.OldPassword, user.PasswordHash); err != nil {
		return errors.New("رمز عبور فعلی اشتباه است")
	}

	newPasswordHash, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = newPasswordHash
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) AdminLogin(ctx context.Context, req *schemas.AdminLoginRequest) (*schemas.AdminAuthLoginResponse, error) {
	admin, err := s.adminRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("username or password is incorrect")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	if !admin.IsActive {
		return nil, errors.New("your account has been deactivated, please contact administrator")
	}

	if err := s.passwordService.CheckPassword(req.Password, admin.PasswordHash); err != nil {
		return nil, errors.New("username or password is incorrect")
	}

	now := time.Now()
	admin.LastLogin = now.Format(time.RFC3339)
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return nil, fmt.Errorf("failed to update last login for admin %s: %w", admin.ID, err)
	}

	return s.generateAdminTokens(ctx, admin)
}

func (s *AuthService) AdminRefreshToken(ctx context.Context, refreshToken string) (*schemas.AdminRefreshTokenResponse, error) {
	admin, err := s.adminRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token expired or invalid")
		}
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	if !admin.IsActive {
		return nil, errors.New("your account has been deactivated, please contact administrator")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(admin.ID.String(), "admin")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, expiry, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	admin.RefreshToken = newRefreshToken
	admin.RefreshTokenExpiry = &expiry
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &schemas.AdminRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.jwtService.GetAccessExpiration().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *AuthService) AdminLogout(ctx context.Context, adminID string) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("admin not found")
		}
		return fmt.Errorf("failed to get admin: %w", err)
	}

	admin.RefreshToken = ""
	admin.RefreshTokenExpiry = nil

	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

func (s *AuthService) AdminChangePassword(ctx context.Context, adminID string, req *schemas.AdminChangePasswordRequest) error {
	admin, err := s.adminRepo.GetByID(ctx, adminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("admin not found")
		}
		return fmt.Errorf("failed to get admin: %w", err)
	}

	if err := s.passwordService.CheckPassword(req.OldPassword, admin.PasswordHash); err != nil {
		return errors.New("current password is incorrect")
	}

	newPasswordHash, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin.PasswordHash = newPasswordHash
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) generateUserTokens(ctx context.Context, user *models.User) (*schemas.UserLoginResponse, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID.String(), "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, expiry, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	user.RefreshToken = refreshToken
	user.RefreshTokenExpiry = &expiry
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save refresh token for user %s: %w", user.ID, err)
	}

	return &schemas.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtService.GetAccessExpiration().Seconds()),
		TokenType:    "Bearer",
		User:         s.toUserResponse(user),
	}, nil
}

func (s *AuthService) generateAdminTokens(ctx context.Context, admin *models.AdminUser) (*schemas.AdminAuthLoginResponse, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(admin.ID.String(), "admin")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, expiry, err := s.jwtService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	admin.RefreshToken = refreshToken
	admin.RefreshTokenExpiry = &expiry
	if err := s.adminRepo.Update(ctx, admin); err != nil {
		return nil, fmt.Errorf("failed to save refresh token for admin %s: %w", admin.ID, err)
	}

	return &schemas.AdminAuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtService.GetAccessExpiration().Seconds()),
		TokenType:    "Bearer",
		AdminUser:    s.toAdminResponse(admin),
	}, nil
}

func (s *AuthService) validatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 11 {
		return errors.New("شماره تماس باید 11 رقم باشد")
	}

	if phoneNumber[0] != '0' || phoneNumber[1] != '9' {
		return errors.New("شماره تماس باید با 09 شروع شود")
	}

	for _, ch := range phoneNumber {
		if ch < '0' || ch > '9' {
			return errors.New("شماره تماس نباید حاوی حرف باشد")
		}
	}

	return nil
}

func (s *AuthService) toUserResponse(user *models.User) schemas.UserResponse {
	return schemas.UserResponse{
		ID:                  user.ID.String(),
		PhoneNumber:         user.PhoneNumber,
		RegistrationDate:    user.RegistrationDate,
		LastLogin:           user.LastLogin,
		LoginCount:          user.LoginCount,
		AgreementAccepted:   user.AgreementAccepted,
		AgreementAcceptedAt: user.AgreementAcceptedAt,
		CloudSyncEnabled:    user.CloudSyncEnabled,
		DoNotDisturbEnabled: user.DoNotDisturbEnabled,
		DNDStartTime:        user.DNDStartTime,
		DNDEndTime:          user.DNDEndTime,
		AndroidVersion:      user.AndroidVersion,
		AppVersion:          user.AppVersion,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func (s *AuthService) toAdminResponse(admin *models.AdminUser) schemas.AdminUserResponse {
	return schemas.AdminUserResponse{
		ID:        admin.ID.String(),
		Username:  admin.Username,
		Email:     admin.Email,
		FullName:  admin.FullName,
		RoleID:    admin.RoleID,
		IsActive:  admin.IsActive,
		LastLogin: admin.LastLogin,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}
