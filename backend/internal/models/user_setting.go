package models

// UserSetting represents user preferences and settings
type UserSetting struct {
	BaseModel
	UserID               string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	NotificationEnabled  bool   `gorm:"default:true" json:"notification_enabled"`
	VibrationEnabled     bool   `gorm:"default:true" json:"vibration_enabled"`
	Language             string `gorm:"type:varchar(5);default:'fa'" json:"language"`
	FontSize             string `gorm:"type:varchar(10);default:'medium'" json:"font_size"`
	Theme                string `gorm:"type:varchar(20);default:'calm'" json:"theme"`
	CrisisAlertThreshold int    `gorm:"default:7" json:"crisis_alert_threshold"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}
