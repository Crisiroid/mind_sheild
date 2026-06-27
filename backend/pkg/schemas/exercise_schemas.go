package schemas

import "time"

type BreathingSessionCreateRequest struct {
	UserID           string `json:"user_id" validate:"required"`
	BreathingPattern string `json:"breathing_pattern,omitempty"`
	DayNumber        *int   `json:"day_number,omitempty"`
}

type BreathingSessionUpdateRequest struct {
	SessionEnd      *time.Time `json:"session_end,omitempty"`
	DurationSeconds *int       `json:"duration_seconds,omitempty"`
	IsCompleted     *bool      `json:"is_completed,omitempty"`
	CalendarTicked  *bool      `json:"calendar_ticked,omitempty"`
}

type BreathingSessionResponse struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	SessionStart     time.Time  `json:"session_start"`
	SessionEnd       *time.Time `json:"session_end,omitempty"`
	DurationSeconds  *int       `json:"duration_seconds,omitempty"`
	BreathingPattern string     `json:"breathing_pattern,omitempty"`
	IsCompleted      bool       `json:"is_completed"`
	CalendarTicked   bool       `json:"calendar_ticked"`
	DayNumber        *int       `json:"day_number,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type BreathingSessionListRequest struct {
	PaginatedRequest
	FilterRequest
	IsCompleted *bool `query:"is_completed" form:"is_completed"`
}

type BreathingSessionListResponse struct {
	PaginatedResponse[BreathingSessionResponse]
}

type CognitiveGameCreateRequest struct {
	UserID           string `json:"user_id" validate:"required"`
	ScenarioID       int    `json:"scenario_id" validate:"required"`
	ScenarioType     string `json:"scenario_type,omitempty"`
	Score            *int   `json:"score,omitempty"`
	IsCorrect        *bool  `json:"is_correct,omitempty"`
	TimeTakenSeconds *int   `json:"time_taken_seconds,omitempty"`
	DayNumber        *int   `json:"day_number,omitempty"`
}

type CognitiveGameResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	GameDate         string    `json:"game_date"`
	ScenarioID       int       `json:"scenario_id"`
	ScenarioType     string    `json:"scenario_type,omitempty"`
	Score            *int      `json:"score,omitempty"`
	IsCorrect        *bool     `json:"is_correct,omitempty"`
	TimeTakenSeconds *int      `json:"time_taken_seconds,omitempty"`
	DayNumber        *int      `json:"day_number,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CognitiveGameUpdateRequest struct {
	Score            *int  `json:"score,omitempty"`
	IsCorrect        *bool `json:"is_correct,omitempty"`
	TimeTakenSeconds *int  `json:"time_taken_seconds,omitempty"`
}

type CognitiveGameListRequest struct {
	PaginatedRequest
	FilterRequest
	ScenarioType *string `query:"scenario_type" form:"scenario_type"`
	IsCorrect    *bool   `query:"is_correct" form:"is_correct"`
}

type CognitiveGameListResponse struct {
	PaginatedResponse[CognitiveGameResponse]
}

type MentalMustCreateRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	MustText  string `json:"must_text" validate:"required"`
	DayNumber *int   `json:"day_number,omitempty"`
}

type MentalMustUpdateRequest struct {
	IsReleased *bool `json:"is_released,omitempty"`
}

type MentalMustResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	MustText     string    `json:"must_text"`
	CreatedDate  string    `json:"created_date"`
	IsReleased   bool      `json:"is_released"`
	ReleasedDate string    `json:"released_date,omitempty"`
	DayNumber    *int      `json:"day_number,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MentalMustListRequest struct {
	PaginatedRequest
	FilterRequest
	IsReleased *bool `query:"is_released" form:"is_released"`
}

type MentalMustListResponse struct {
	PaginatedResponse[MentalMustResponse]
}

type NegativeThoughtCreateRequest struct {
	UserID             string `json:"user_id" validate:"required"`
	ThoughtText        string `json:"thought_text" validate:"required"`
	Situation          string `json:"situation,omitempty"`
	CognitiveErrorType string `json:"cognitive_error_type,omitempty"`
	ImpactLevel        *int   `json:"impact_level,omitempty" validate:"omitempty,min=1,max=10"`
	DayNumber          *int   `json:"day_number,omitempty"`
}

type NegativeThoughtUpdateRequest struct {
	Situation          *string `json:"situation,omitempty"`
	CognitiveErrorType *string `json:"cognitive_error_type,omitempty"`
	ImpactLevel        *int    `json:"impact_level,omitempty" validate:"omitempty,min=1,max=10"`
}

type NegativeThoughtResponse struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	ThoughtText        string    `json:"thought_text"`
	Situation          string    `json:"situation,omitempty"`
	CognitiveErrorType string    `json:"cognitive_error_type,omitempty"`
	ImpactLevel        *int      `json:"impact_level,omitempty"`
	RecordedAt         string    `json:"recorded_at"`
	DayNumber          *int      `json:"day_number,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type NegativeThoughtListRequest struct {
	PaginatedRequest
	FilterRequest
	CognitiveErrorType *string `query:"cognitive_error_type" form:"cognitive_error_type"`
	ImpactMin          *int    `query:"impact_min" form:"impact_min"`
	ImpactMax          *int    `query:"impact_max" form:"impact_max"`
}

type NegativeThoughtListResponse struct {
	PaginatedResponse[NegativeThoughtResponse]
}

type MindCourtCreateRequest struct {
	UserID                string `json:"user_id" validate:"required"`
	NegativeThoughtID     string `json:"negative_thought_id" validate:"required"`
	SupportingEvidence    string `json:"supporting_evidence,omitempty"`
	ContradictingEvidence string `json:"contradicting_evidence,omitempty"`
	GuideHelperUsed       bool   `json:"guide_helper_used,omitempty"`
	AlternativeThought    string `json:"alternative_thought,omitempty"`
	DayNumber             *int   `json:"day_number,omitempty"`
}

type MindCourtUpdateRequest struct {
	SupportingEvidence    *string `json:"supporting_evidence,omitempty"`
	ContradictingEvidence *string `json:"contradicting_evidence,omitempty"`
	AlternativeThought    *string `json:"alternative_thought,omitempty"`
}

type MindCourtResponse struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	NegativeThoughtID     string    `json:"negative_thought_id"`
	SupportingEvidence    string    `json:"supporting_evidence,omitempty"`
	ContradictingEvidence string    `json:"contradicting_evidence,omitempty"`
	GuideHelperUsed       bool      `json:"guide_helper_used"`
	AlternativeThought    string    `json:"alternative_thought,omitempty"`
	CreatedDate           string    `json:"created_date"`
	DayNumber             *int      `json:"day_number,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type MindCourtListRequest struct {
	PaginatedRequest
	FilterRequest
	NegativeThoughtID *string `query:"negative_thought_id" form:"negative_thought_id"`
	GuideHelperUsed   *bool   `query:"guide_helper_used" form:"guide_helper_used"`
}

type MindCourtListResponse struct {
	PaginatedResponse[MindCourtResponse]
}

type ConflictExerciseCreateRequest struct {
	UserID           string `json:"user_id" validate:"required"`
	ScenarioID       int    `json:"scenario_id" validate:"required"`
	PerformanceScore *int   `json:"performance_score,omitempty"`
	DayNumber        *int   `json:"day_number,omitempty"`
}

type ConflictExerciseUpdateRequest struct {
	PracticeCount    *int `json:"practice_count,omitempty"`
	PerformanceScore *int `json:"performance_score,omitempty"`
}

type ConflictExerciseResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	ScenarioID       int       `json:"scenario_id"`
	PracticeCount    int       `json:"practice_count"`
	LastPracticeDate string    `json:"last_practice_date,omitempty"`
	PerformanceScore *int      `json:"performance_score,omitempty"`
	DayNumber        *int      `json:"day_number,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ConflictExerciseListRequest struct {
	PaginatedRequest
	FilterRequest
	ScenarioID *int `query:"scenario_id" form:"scenario_id"`
}

type ConflictExerciseListResponse struct {
	PaginatedResponse[ConflictExerciseResponse]
}

type MoodTrackerCreateRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	ActivityID   string `json:"activity_id,omitempty"`
	ActivityName string `json:"activity_name,omitempty"`
	MoodBefore   int    `json:"mood_before" validate:"required,min=1,max=10"`
	MoodAfter    int    `json:"mood_after" validate:"required,min=1,max=10"`
	Notes        string `json:"notes,omitempty"`
	DayNumber    *int   `json:"day_number,omitempty"`
}

type MoodTrackerUpdateRequest struct {
	MoodAfter *int    `json:"mood_after,omitempty" validate:"omitempty,min=1,max=10"`
	Notes     *string `json:"notes,omitempty"`
}

type MoodTrackerResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ActivityID   string    `json:"activity_id,omitempty"`
	ActivityName string    `json:"activity_name,omitempty"`
	MoodBefore   int       `json:"mood_before"`
	MoodAfter    int       `json:"mood_after"`
	MoodDelta    int       `json:"mood_delta"`
	ActivityDate string    `json:"activity_date"`
	Notes        string    `json:"notes,omitempty"`
	DayNumber    *int      `json:"day_number,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MoodTrackerListRequest struct {
	PaginatedRequest
	FilterRequest
	ActivityName *string `query:"activity_name" form:"activity_name"`
	MoodMin      *int    `query:"mood_min" form:"mood_min"`
	MoodMax      *int    `query:"mood_max" form:"mood_max"`
}

type MoodTrackerListResponse struct {
	PaginatedResponse[MoodTrackerResponse]
}

type RoleValueCreateRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	EntryType string `json:"entry_type" validate:"required,oneof=role value"`
	EntryText string `json:"entry_text" validate:"required"`
	DayNumber *int   `json:"day_number,omitempty"`
}

type RoleValueUpdateRequest struct {
	EntryText *string `json:"entry_text,omitempty"`
}

type RoleValueResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	EntryType   string    `json:"entry_type"`
	EntryText   string    `json:"entry_text"`
	CreatedDate string    `json:"created_date"`
	DayNumber   *int      `json:"day_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleValueListRequest struct {
	PaginatedRequest
	FilterRequest
	EntryType *string `query:"entry_type" form:"entry_type"`
}

type RoleValueListResponse struct {
	PaginatedResponse[RoleValueResponse]
}

type SkyThoughtCreateRequest struct {
	UserID      string `json:"user_id" validate:"required"`
	ThoughtText string `json:"thought_text" validate:"required"`
	DayNumber   *int   `json:"day_number,omitempty"`
}

type SkyThoughtUpdateRequest struct {
	CloudSwiped *bool `json:"cloud_swiped,omitempty"`
}

type SkyThoughtResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ThoughtText string    `json:"thought_text"`
	CloudSwiped bool      `json:"cloud_swiped"`
	SwipedAt    string    `json:"swiped_at,omitempty"`
	CreatedDate string    `json:"created_date"`
	DayNumber   *int      `json:"day_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SkyThoughtListRequest struct {
	PaginatedRequest
	FilterRequest
	CloudSwiped *bool `query:"cloud_swiped" form:"cloud_swiped"`
}

type SkyThoughtListResponse struct {
	PaginatedResponse[SkyThoughtResponse]
}
