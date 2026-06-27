package models

// ConflictExercise represents work conflict scenario exercises
type ConflictExercise struct {
	BaseModel
	UserID           string `gorm:"type:uuid;not null" json:"user_id"`
	ScenarioID       int    `gorm:"uniqueIndex:idx_user_scenario;not null" json:"scenario_id"`
	PracticeCount    int    `gorm:"default:1" json:"practice_count"`
	LastPracticeDate string `json:"last_practice_date,omitempty"`
	PerformanceScore *int   `json:"performance_score,omitempty"`
	DayNumber        *int   `json:"day_number,omitempty"`
}

func (ConflictExercise) TableName() string {
	return "conflict_exercises"
}
