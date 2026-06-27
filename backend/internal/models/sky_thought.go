package models

// SkyThought represents thought sky cloud entries
type SkyThought struct {
	BaseModel
	UserID      string `gorm:"type:uuid;not null" json:"user_id"`
	ThoughtText string `gorm:"type:text;not null" json:"thought_text"`
	CloudSwiped bool   `gorm:"default:false" json:"cloud_swiped"`
	SwipedAt    string `json:"swiped_at,omitempty"`
	CreatedDate string `gorm:"autoCreateTime" json:"created_date"`
	DayNumber   *int   `json:"day_number,omitempty"`
}

func (SkyThought) TableName() string {
	return "sky_thoughts"
}
