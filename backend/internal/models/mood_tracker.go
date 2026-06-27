package models

// MoodTracker represents mood and activity tracking
type MoodTracker struct {
	BaseModel
	UserID       string `gorm:"type:uuid;not null" json:"user_id"`
	ActivityID   string `gorm:"type:uuid" json:"activity_id,omitempty"`
	ActivityName string `gorm:"type:varchar(100)" json:"activity_name,omitempty"`
	MoodBefore   int    `gorm:"check:mood_before BETWEEN 1 AND 10;not null" json:"mood_before"`
	MoodAfter    int    `gorm:"check:mood_after BETWEEN 1 AND 10;not null" json:"mood_after"`
	ActivityDate string `gorm:"autoCreateTime" json:"activity_date"`
	Notes        string `gorm:"type:text" json:"notes,omitempty"`
	DayNumber    *int   `json:"day_number,omitempty"`
}

func (MoodTracker) TableName() string {
	return "mood_tracker"
}
