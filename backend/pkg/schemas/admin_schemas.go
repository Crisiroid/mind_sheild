package schemas

import "time"

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Token     string            `json:"token"`
	AdminUser AdminUserResponse `json:"admin_user"`
}

type AdminCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name,omitempty"`
	RoleID   string `json:"role_id,omitempty"`
	IsActive bool   `json:"is_active"`
}

type AdminUpdateRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	FullName *string `json:"full_name,omitempty"`
	RoleID   *string `json:"role_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type AdminUserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name,omitempty"`
	RoleID    string    `json:"role_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	LastLogin string    `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminListRequest struct {
	PaginatedRequest
	FilterRequest
	RoleID   *string `query:"role_id" form:"role_id"`
	IsActive *bool   `query:"is_active" form:"is_active"`
}

type AdminListResponse struct {
	PaginatedResponse[AdminUserResponse]
}

type AdminRoleCreateRequest struct {
	RoleName    string `json:"role_name" validate:"required"`
	Description string `json:"description,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

type AdminRoleUpdateRequest struct {
	Description *string `json:"description,omitempty"`
	Permissions *string `json:"permissions,omitempty"`
}

type AdminRoleResponse struct {
	ID          string    `json:"id"`
	RoleName    string    `json:"role_name"`
	Description string    `json:"description,omitempty"`
	Permissions string    `json:"permissions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminRoleListRequest struct {
	PaginatedRequest
	FilterRequest
}

type AdminRoleListResponse struct {
	PaginatedResponse[AdminRoleResponse]
}

type UserReportCreateRequest struct {
	ReportType         string  `json:"report_type" validate:"required"`
	ReportDate         string  `json:"report_date" validate:"required"`
	TotalUsers         int     `json:"total_users"`
	ActiveUsers        int     `json:"active_users"`
	AvgEngagementScore float64 `json:"avg_engagement_score,omitempty"`
	CrisisAlertsCount  int     `json:"crisis_alerts_count"`
	AnonymizedData     string  `json:"anonymized_data,omitempty"`
}

type UserReportResponse struct {
	ID                 string    `json:"id"`
	ReportType         string    `json:"report_type"`
	ReportDate         string    `json:"report_date"`
	TotalUsers         int       `json:"total_users"`
	ActiveUsers        int       `json:"active_users"`
	AvgEngagementScore float64   `json:"avg_engagement_score,omitempty"`
	CrisisAlertsCount  int       `json:"crisis_alerts_count"`
	AnonymizedData     string    `json:"anonymized_data,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type UserReportListRequest struct {
	PaginatedRequest
	FilterRequest
	ReportType *string `query:"report_type" form:"report_type"`
}

type UserReportListResponse struct {
	PaginatedResponse[UserReportResponse]
}

type SystemLogResponse struct {
	ID         string    `json:"id"`
	LogType    string    `json:"log_type"`
	LogMessage string    `json:"log_message,omitempty"`
	UserID     string    `json:"user_id,omitempty"`
	Severity   string    `json:"severity"`
	CreatedAt  time.Time `json:"created_at"`
}

type SystemLogListRequest struct {
	PaginatedRequest
	FilterRequest
	LogType  *string `query:"log_type" form:"log_type"`
	Severity *string `query:"severity" form:"severity"`
}

type SystemLogListResponse struct {
	PaginatedResponse[SystemLogResponse]
}
