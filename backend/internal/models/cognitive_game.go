package models

type CognitiveErrorGame struct {
	BaseModel
	UserID           string `gorm:"type:uuid;not null" json:"user_id"`
	GameDate         string `gorm:"autoCreateTime" json:"game_date"`
	ScenarioID       int    `gorm:"not null" json:"scenario_id"`
	ScenarioType     string `gorm:"type:varchar(50)" json:"scenario_type,omitempty"`
	Score            *int   `json:"score,omitempty"`
	IsCorrect        *bool  `json:"is_correct,omitempty"`
	TimeTakenSeconds *int   `json:"time_taken_seconds,omitempty"`
	DayNumber        *int   `json:"day_number,omitempty"`
}

func (CognitiveErrorGame) TableName() string {
	return "cognitive_error_games"
}
