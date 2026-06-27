package models

type NegativeThought struct {
	BaseModel
	UserID             string `gorm:"type:uuid;not null" json:"user_id"`
	ThoughtText        string `gorm:"type:text;not null" json:"thought_text"`
	Situation          string `gorm:"type:text" json:"situation,omitempty"`
	CognitiveErrorType string `gorm:"type:varchar(50)" json:"cognitive_error_type,omitempty"`
	ImpactLevel        *int   `gorm:"check:impact_level BETWEEN 1 AND 10" json:"impact_level,omitempty"`
	RecordedAt         string `gorm:"autoCreateTime" json:"recorded_at"`
	DayNumber          *int   `json:"day_number,omitempty"`
}

func (NegativeThought) TableName() string {
	return "negative_thoughts"
}
