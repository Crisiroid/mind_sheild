package models

import "time"

type BreathingSession struct {
	BaseModel
	UserID           string     `gorm:"type:uuid;not null" json:"user_id"`
	SessionStart     time.Time  `gorm:"autoCreateTime" json:"session_start"`
	SessionEnd       *time.Time `json:"session_end,omitempty"`
	DurationSeconds  *int       `json:"duration_seconds,omitempty"`
	BreathingPattern string     `gorm:"type:varchar(50)" json:"breathing_pattern,omitempty"`
	IsCompleted      bool       `gorm:"default:false" json:"is_completed"`
	CalendarTicked   bool       `gorm:"default:false" json:"calendar_ticked"`
	DayNumber        *int       `json:"day_number,omitempty"`
}

func (BreathingSession) TableName() string {
	return "breathing_sessions"
}
