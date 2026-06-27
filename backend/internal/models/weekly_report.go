package models

// WeeklyReport represents weekly progress reports
type WeeklyReport struct {
	BaseModel
	UserID                 string  `gorm:"type:uuid;not null" json:"user_id"`
	WeekNumber             int     `gorm:"uniqueIndex:idx_user_week;check:week_number BETWEEN 1 AND 8;not null" json:"week_number"`
	StartDate              string  `gorm:"not null" json:"start_date"`
	EndDate                string  `gorm:"not null" json:"end_date"`
	StressEventsCount      int     `gorm:"default:0" json:"stress_events_count"`
	AvgStressIntensity     float64 `json:"avg_stress_intensity,omitempty"`
	BreathingSessionsCount int     `gorm:"default:0" json:"breathing_sessions_count"`
	NegativeThoughtsCount  int     `gorm:"default:0" json:"negative_thoughts_count"`
	BodyTensionMapsCount   int     `gorm:"default:0" json:"body_tension_maps_count"`
	MoodImprovementScore   float64 `json:"mood_improvement_score,omitempty"`
	ActivitiesDistribution string  `gorm:"type:jsonb" json:"activities_distribution,omitempty"`
	ProgressPercentage     float64 `json:"progress_percentage,omitempty"`
	GeneratedAt            string  `gorm:"autoCreateTime" json:"generated_at"`
}

func (WeeklyReport) TableName() string {
	return "weekly_reports"
}
