package models

// SystemLog represents system activity logs
type SystemLog struct {
	BaseModel
	LogType    string `gorm:"type:varchar(50);not null" json:"log_type"`
	LogMessage string `gorm:"type:text" json:"log_message,omitempty"`
	UserID     string `gorm:"type:uuid" json:"user_id,omitempty"`
	Severity   string `gorm:"type:varchar(20);default:'info'" json:"severity"`
}

func (SystemLog) TableName() string {
	return "admin_panel.system_logs"
}
