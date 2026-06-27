package models

type MentalMust struct {
	BaseModel
	UserID       string `gorm:"type:uuid;not null" json:"user_id"`
	MustText     string `gorm:"type:text;not null" json:"must_text"`
	CreatedDate  string `gorm:"autoCreateTime" json:"created_date"`
	IsReleased   bool   `gorm:"default:false" json:"is_released"`
	ReleasedDate string `json:"released_date,omitempty"`
	DayNumber    *int   `json:"day_number,omitempty"`
}

func (MentalMust) TableName() string {
	return "mental_musts"
}
