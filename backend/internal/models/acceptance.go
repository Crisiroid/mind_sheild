package models

// AcceptanceExercise represents active acceptance vs surrender exercises
type AcceptanceExercise struct {
	BaseModel
	UserID             string `gorm:"type:uuid;not null" json:"user_id"`
	VideoWatched       bool   `gorm:"default:false" json:"video_watched"`
	WatchedAt          string `json:"watched_at,omitempty"`
	UnderstandingLevel *int   `gorm:"check:understanding_level BETWEEN 1 AND 10" json:"understanding_level,omitempty"`
	Notes              string `gorm:"type:text" json:"notes,omitempty"`
	DayNumber          *int   `json:"day_number,omitempty"`
}

func (AcceptanceExercise) TableName() string {
	return "acceptance_exercises"
}
