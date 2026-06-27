package models

import "time"

type BodyTensionMap struct {
	BaseModel
	UserID           string    `gorm:"type:uuid;not null" json:"user_id"`
	MappingDate      time.Time `gorm:"autoCreateTime" json:"mapping_date"`
	BodyRegions      string    `gorm:"type:jsonb;not null" json:"body_regions"`
	OverallIntensity *int      `gorm:"check:overall_intensity BETWEEN 1 AND 10" json:"overall_intensity,omitempty"`
	SeverityColor    string    `gorm:"type:varchar(20)" json:"severity_color,omitempty"`
	Notes            string    `gorm:"type:text" json:"notes,omitempty"`
	DayNumber        *int      `json:"day_number,omitempty"`
}

func (BodyTensionMap) TableName() string {
	return "body_tension_maps"
}
