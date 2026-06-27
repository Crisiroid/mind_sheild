package models

import "time"

// User represents a user account in the psychology app
type User struct {
	BaseModel
	PhoneNumber         string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone_number"`
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

	// Relationships
	Settings UserSetting `gorm:"foreignKey:UserID" json:"settings,omitempty"`
}

func (User) TableName() string {
	return "users"
}
