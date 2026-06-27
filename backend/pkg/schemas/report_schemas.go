package schemas

import "time"

type BaseStats struct {
	Total   int64   `json:"total"`
	Average float64 `json:"average,omitempty"`
	Min     float64 `json:"min,omitempty"`
	Max     float64 `json:"max,omitempty"`
	Sum     float64 `json:"sum,omitempty"`
}

type DistributionStats struct {
	Label string  `json:"label"`
	Count int64   `json:"count"`
	Value float64 `json:"value,omitempty"`
}

type TrendDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
	Value     float64   `json:"value,omitempty"`
	Label     string    `json:"label,omitempty"`
}

type UserStatsResponse struct {
	TotalUsers        int64   `json:"total_users"`
	ActiveUsers       int64   `json:"active_users"`
	NewUsersToday     int64   `json:"new_users_today"`
	NewUsersThisWeek  int64   `json:"new_users_this_week"`
	NewUsersThisMonth int64   `json:"new_users_this_month"`
	AgreementRate     float64 `json:"agreement_rate"`
	AvgLoginCount     float64 `json:"avg_login_count"`
}

type ActivityStatsResponse struct {
	TotalActivities int64   `json:"total_activities"`
	CompletedCount  int64   `json:"completed_count"`
	IncompleteCount int64   `json:"incomplete_count"`
	CompletionRate  float64 `json:"completion_rate"`
	AvgPerUser      float64 `json:"avg_per_user"`
	DateRange       string  `json:"date_range"`
}

type EngagementStatsResponse struct {
	TotalUsers       int64   `json:"total_users"`
	ActiveUsers      int64   `json:"active_users"`
	EngagementRate   float64 `json:"engagement_rate"`
	AvgActivities    float64 `json:"avg_activities"`
	HighEngagement   int64   `json:"high_engagement"`
	MediumEngagement int64   `json:"medium_engagement"`
	LowEngagement    int64   `json:"low_engagement"`
}

type ExportRecord struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id,omitempty"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type ReportFilter struct {
	DateFrom  *time.Time `json:"date_from"`
	DateTo    *time.Time `json:"date_to"`
	UserID    string     `json:"user_id"`
	DayNumber *int       `json:"day_number"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
	SortBy    string     `json:"sort_by"`
	SortOrder string     `json:"sort_order"`
	Search    string     `json:"search"`
}

type CompletionStatsResponse struct {
	Total          int64   `json:"total"`
	Completed      int64   `json:"completed"`
	Incomplete     int64   `json:"incomplete"`
	CompletionRate float64 `json:"completion_rate"`
	AvgCompletion  float64 `json:"avg_completion"`
}

type PerformanceStatsResponse struct {
	TotalEntries    int64   `json:"total_entries"`
	AvgScore        float64 `json:"avg_score"`
	MinScore        float64 `json:"min_score"`
	MaxScore        float64 `json:"max_score"`
	ImprovementRate float64 `json:"improvement_rate"`
}

type TimeAnalysisResponse struct {
	AvgDuration   float64 `json:"avg_duration"`
	MinDuration   float64 `json:"min_duration"`
	MaxDuration   float64 `json:"max_duration"`
	TotalDuration float64 `json:"total_duration"`
	EntriesCount  int64   `json:"entries_count"`
}

type IntensityStatsResponse struct {
	AvgIntensity    float64 `json:"avg_intensity"`
	MinIntensity    float64 `json:"min_intensity"`
	MaxIntensity    float64 `json:"max_intensity"`
	TotalEntries    int64   `json:"total_entries"`
	HighIntensity   int64   `json:"high_intensity"`
	MediumIntensity int64   `json:"medium_intensity"`
	LowIntensity    int64   `json:"low_intensity"`
}

type WeeklyReportCreateRequest struct {
	UserID                 string  `json:"user_id" validate:"required"`
	WeekNumber             int     `json:"week_number" validate:"required,min=1,max=8"`
	StartDate              string  `json:"start_date" validate:"required"`
	EndDate                string  `json:"end_date" validate:"required"`
	StressEventsCount      int     `json:"stress_events_count"`
	AvgStressIntensity     float64 `json:"avg_stress_intensity,omitempty"`
	BreathingSessionsCount int     `json:"breathing_sessions_count"`
	NegativeThoughtsCount  int     `json:"negative_thoughts_count"`
	BodyTensionMapsCount   int     `json:"body_tension_maps_count"`
	MoodImprovementScore   float64 `json:"mood_improvement_score,omitempty"`
	ActivitiesDistribution string  `json:"activities_distribution,omitempty"`
	ProgressPercentage     float64 `json:"progress_percentage,omitempty"`
}

type WeeklyReportUpdateRequest struct {
	StressEventsCount      *int     `json:"stress_events_count,omitempty"`
	AvgStressIntensity     *float64 `json:"avg_stress_intensity,omitempty"`
	BreathingSessionsCount *int     `json:"breathing_sessions_count,omitempty"`
	NegativeThoughtsCount  *int     `json:"negative_thoughts_count,omitempty"`
	BodyTensionMapsCount   *int     `json:"body_tension_maps_count,omitempty"`
	MoodImprovementScore   *float64 `json:"mood_improvement_score,omitempty"`
	ProgressPercentage     *float64 `json:"progress_percentage,omitempty"`
}

type WeeklyReportResponse struct {
	ID                     string    `json:"id"`
	UserID                 string    `json:"user_id"`
	WeekNumber             int       `json:"week_number"`
	StartDate              string    `json:"start_date"`
	EndDate                string    `json:"end_date"`
	StressEventsCount      int       `json:"stress_events_count"`
	AvgStressIntensity     float64   `json:"avg_stress_intensity,omitempty"`
	BreathingSessionsCount int       `json:"breathing_sessions_count"`
	NegativeThoughtsCount  int       `json:"negative_thoughts_count"`
	BodyTensionMapsCount   int       `json:"body_tension_maps_count"`
	MoodImprovementScore   float64   `json:"mood_improvement_score,omitempty"`
	ActivitiesDistribution string    `json:"activities_distribution,omitempty"`
	ProgressPercentage     float64   `json:"progress_percentage,omitempty"`
	GeneratedAt            string    `json:"generated_at"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type WeeklyReportListRequest struct {
	PaginatedRequest
	FilterRequest
	WeekNumber *int `query:"week_number" form:"week_number"`
}

type WeeklyReportListResponse struct {
	PaginatedResponse[WeeklyReportResponse]
}
