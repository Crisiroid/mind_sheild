package models

import "time"

type StressEvent struct {
	BaseModel
	UserID               string    `gorm:"type:uuid;not null" json:"user_id"`
	EventTimestamp       time.Time `gorm:"autoCreateTime" json:"event_timestamp"`
	SituationType        string    `gorm:"type:varchar(50);not null" json:"situation_type"`
	SituationDescription string    `gorm:"type:text" json:"situation_description,omitempty"`
	IntensityLevel       int       `gorm:"check:intensity_level BETWEEN 1 AND 10;not null" json:"intensity_level"`
	Location             string    `gorm:"type:varchar(100)" json:"location,omitempty"`
	DayNumber            *int      `json:"day_number,omitempty"`
}

func (StressEvent) TableName() string {
	return "stress_events"
}
