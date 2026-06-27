package schemas

type UserRegisterRequest struct {
	PhoneNumber    string `json:"phone_number" validate:"required"`
	Password       string `json:"password" validate:"required,min=6"`
	AndroidVersion string `json:"android_version,omitempty"`
	AppVersion     string `json:"app_version,omitempty"`
}

type UserLoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}

type UserRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type AdminAuthLoginResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int64             `json:"expires_in"`
	TokenType    string            `json:"token_type"`
	AdminUser    AdminUserResponse `json:"admin_user"`
}

type AdminRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AdminRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type AdminChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type TokenClaims struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
	Email    string `json:"email,omitempty"`
	IssuedAt int64  `json:"issued_at"`
}
