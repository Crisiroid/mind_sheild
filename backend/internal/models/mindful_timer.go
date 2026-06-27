package models

type MindfulTimer struct {
	BaseModel
	UserID                  string `gorm:"type:uuid;not null" json:"user_id"`
	TimerStart              string `gorm:"autoCreateTime" json:"timer_start"`
	TimerEnd                string `json:"timer_end,omitempty"`
	DurationSeconds         *int   `json:"duration_seconds,omitempty"`
	VibrationRemindersCount int    `gorm:"default:0" json:"vibration_reminders_count"`
	IsCompleted             bool   `gorm:"default:false" json:"is_completed"`
	DayNumber               *int   `json:"day_number,omitempty"`
}

func (MindfulTimer) TableName() string {
	return "mindful_timers"
}
