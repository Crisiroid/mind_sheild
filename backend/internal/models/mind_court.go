package models

type MindCourtEvidence struct {
	BaseModel
	UserID                string `gorm:"type:uuid;not null" json:"user_id"`
	NegativeThoughtID     string `gorm:"type:uuid" json:"negative_thought_id"`
	SupportingEvidence    string `gorm:"type:text" json:"supporting_evidence,omitempty"`
	ContradictingEvidence string `gorm:"type:text" json:"contradicting_evidence,omitempty"`
	GuideHelperUsed       bool   `gorm:"default:false" json:"guide_helper_used"`
	AlternativeThought    string `gorm:"type:text" json:"alternative_thought,omitempty"`
	CreatedDate           string `gorm:"autoCreateTime" json:"created_date"`
	DayNumber             *int   `json:"day_number,omitempty"`
}

func (MindCourtEvidence) TableName() string {
	return "mind_court_evidence"
}
