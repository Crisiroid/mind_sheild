package models

import "time"

type EmotionTriangleInteraction struct {
	BaseModel
	UserID                string    `gorm:"type:uuid;not null" json:"user_id"`
	InteractionDate       time.Time `gorm:"autoCreateTime" json:"interaction_date"`
	SideClicked           string    `gorm:"type:varchar(20);check:side_clicked IN ('thought', 'body', 'behavior');not null" json:"side_clicked"`
	ThoughtAccountsViewed string    `gorm:"type:jsonb" json:"thought_accounts_viewed,omitempty"`
	VibrationTriggered    bool      `gorm:"default:false" json:"vibration_triggered"`
	DayNumber             *int      `json:"day_number,omitempty"`
}

func (EmotionTriangleInteraction) TableName() string {
	return "emotion_triangle_interactions"
}
