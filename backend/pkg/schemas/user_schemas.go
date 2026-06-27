package schemas

import "time"

type UserCreateRequest struct {
	PhoneNumber    string `json:"phone_number" validate:"required"`
	AndroidVersion string `json:"android_version,omitempty"`
	AppVersion     string `json:"app_version,omitempty"`
}

type UserUpdateRequest struct {
	CloudSyncEnabled    *bool      `json:"cloud_sync_enabled,omitempty"`
	DoNotDisturbEnabled *bool      `json:"do_not_disturb_enabled,omitempty"`
	DNDStartTime        *time.Time `json:"dnd_start_time,omitempty"`
	DNDEndTime          *time.Time `json:"dnd_end_time,omitempty"`
}

type UserResponse struct {
	ID                  string               `json:"id"`
	PhoneNumber         string               `json:"phone_number"`
	RegistrationDate    time.Time            `json:"registration_date"`
	LastLogin           *time.Time           `json:"last_login,omitempty"`
	LoginCount          int                  `json:"login_count"`
	AgreementAccepted   bool                 `json:"agreement_accepted"`
	AgreementAcceptedAt *time.Time           `json:"agreement_accepted_at,omitempty"`
	CloudSyncEnabled    bool                 `json:"cloud_sync_enabled"`
	DoNotDisturbEnabled bool                 `json:"do_not_disturb_enabled"`
	DNDStartTime        *time.Time           `json:"dnd_start_time,omitempty"`
	DNDEndTime          *time.Time           `json:"dnd_end_time,omitempty"`
	AndroidVersion      string               `json:"android_version,omitempty"`
	AppVersion          string               `json:"app_version,omitempty"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
	Settings            *UserSettingResponse `json:"settings,omitempty"`
}

type UserListRequest struct {
	PaginatedRequest
	FilterRequest
	AgreementAccepted *bool `query:"agreement_accepted" form:"agreement_accepted"`
}

type UserListResponse struct {
	PaginatedResponse[UserResponse]
}

type CalendarCreateRequest struct {
	UserID              string   `json:"user_id" validate:"required"`
	DayNumber           int      `json:"day_number" validate:"required,min=1,max=56"`
	CalendarDate        string   `json:"calendar_date" validate:"required"`
	ActivitiesCompleted []string `json:"activities_completed,omitempty"`
}

type CalendarUpdateRequest struct {
	IsCompleted         *bool    `json:"is_completed,omitempty"`
	ActivitiesCompleted []string `json:"activities_completed,omitempty"`
}

type CalendarResponse struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	DayNumber           int        `json:"day_number"`
	CalendarDate        string     `json:"calendar_date"`
	IsCompleted         bool       `json:"is_completed"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	ActivitiesCompleted []string   `json:"activities_completed"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CalendarListRequest struct {
	PaginatedRequest
	FilterRequest
	IsCompleted *bool `query:"is_completed" form:"is_completed"`
}

type CalendarListResponse struct {
	PaginatedResponse[CalendarResponse]
}

type UserSettingUpdateRequest struct {
	NotificationEnabled  *bool   `json:"notification_enabled,omitempty"`
	VibrationEnabled     *bool   `json:"vibration_enabled,omitempty"`
	Language             *string `json:"language,omitempty"`
	FontSize             *string `json:"font_size,omitempty"`
	Theme                *string `json:"theme,omitempty"`
	CrisisAlertThreshold *int    `json:"crisis_alert_threshold,omitempty"`
}

type UserSettingResponse struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	NotificationEnabled  bool      `json:"notification_enabled"`
	VibrationEnabled     bool      `json:"vibration_enabled"`
	Language             string    `json:"language"`
	FontSize             string    `json:"font_size"`
	Theme                string    `json:"theme"`
	CrisisAlertThreshold int       `json:"crisis_alert_threshold"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type UserSettingListRequest struct {
	PaginatedRequest
	FilterRequest
	Language *string `query:"language" form:"language"`
	Theme    *string `query:"theme" form:"theme"`
}

type UserSettingListResponse struct {
	PaginatedResponse[UserSettingResponse]
}
