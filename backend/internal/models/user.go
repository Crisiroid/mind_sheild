package models

import "time"

type User struct {
	BaseModel
	PhoneNumber         string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	PasswordHash        string     `gorm:"type:varchar(255);not null" json:"-"`
	RegistrationDate    time.Time  `gorm:"autoCreateTime" json:"registration_date"`
	LastLogin           *time.Time `json:"last_login,omitempty"`
	LoginCount          int        `gorm:"default:1" json:"login_count"`
	AgreementAccepted   bool       `gorm:"default:false" json:"agreement_accepted"`
	AgreementAcceptedAt *time.Time `json:"agreement_accepted_at,omitempty"`
	CloudSyncEnabled    bool       `gorm:"default:false" json:"cloud_sync_enabled"`
	DoNotDisturbEnabled bool       `gorm:"default:false" json:"do_not_disturb_enabled"`
	DNDStartTime        *time.Time `json:"dnd_start_time,omitempty"`
	DNDEndTime          *time.Time `json:"dnd_end_time,omitempty"`
	AndroidVersion      string     `gorm:"type:varchar(10)" json:"android_version,omitempty"`
	AppVersion          string     `gorm:"type:varchar(10)" json:"app_version,omitempty"`

	RefreshToken       string     `gorm:"type:varchar(500)" json:"-"`
	RefreshTokenExpiry *time.Time `json:"-"`
}

func (User) TableName() string {
	return "users"
}
