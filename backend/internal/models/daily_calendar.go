package models

import "time"

type DailyCalendar struct {
	BaseModel
	UserID              string     `gorm:"type:uuid;uniqueIndex:idx_user_day;not null" json:"user_id"`
	DayNumber           int        `gorm:"uniqueIndex:idx_user_day;check:day_number BETWEEN 1 AND 56" json:"day_number"`
	CalendarDate        time.Time  `gorm:"not null" json:"calendar_date"`
	IsCompleted         bool       `gorm:"default:false" json:"is_completed"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	ActivitiesCompleted string     `gorm:"type:jsonb;default:'[]'" json:"activities_completed"`
}

func (DailyCalendar) TableName() string {
	return "daily_calendar"
}
