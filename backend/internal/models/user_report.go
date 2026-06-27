package models

// UserReport represents anonymized user reports
type UserReport struct {
	BaseModel
	ReportType         string  `gorm:"type:varchar(50);not null" json:"report_type"`
	ReportDate         string  `gorm:"not null" json:"report_date"`
	TotalUsers         int     `gorm:"default:0" json:"total_users"`
	ActiveUsers        int     `gorm:"default:0" json:"active_users"`
	AvgEngagementScore float64 `json:"avg_engagement_score,omitempty"`
	CrisisAlertsCount  int     `gorm:"default:0" json:"crisis_alerts_count"`
	AnonymizedData     string  `gorm:"type:jsonb" json:"anonymized_data,omitempty"`
}

func (UserReport) TableName() string {
	return "admin_panel.user_reports"
}
