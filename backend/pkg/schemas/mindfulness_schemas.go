package schemas

import "time"

type MindfulTimerCreateRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	DayNumber *int   `json:"day_number,omitempty"`
}

type MindfulTimerUpdateRequest struct {
	TimerEnd                *string `json:"timer_end,omitempty"`
	DurationSeconds         *int    `json:"duration_seconds,omitempty"`
	VibrationRemindersCount *int    `json:"vibration_reminders_count,omitempty"`
	IsCompleted             *bool   `json:"is_completed,omitempty"`
}

type MindfulTimerResponse struct {
	ID                      string    `json:"id"`
	UserID                  string    `json:"user_id"`
	TimerStart              string    `json:"timer_start"`
	TimerEnd                string    `json:"timer_end,omitempty"`
	DurationSeconds         *int      `json:"duration_seconds,omitempty"`
	VibrationRemindersCount int       `json:"vibration_reminders_count"`
	IsCompleted             bool      `json:"is_completed"`
	DayNumber               *int      `json:"day_number,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type MindfulTimerListRequest struct {
	PaginatedRequest
	FilterRequest
	IsCompleted *bool `query:"is_completed" form:"is_completed"`
}

type MindfulTimerListResponse struct {
	PaginatedResponse[MindfulTimerResponse]
}

type AcceptanceExerciseCreateRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	DayNumber *int   `json:"day_number,omitempty"`
}

type AcceptanceExerciseUpdateRequest struct {
	VideoWatched       *bool   `json:"video_watched,omitempty"`
	UnderstandingLevel *int    `json:"understanding_level,omitempty" validate:"omitempty,min=1,max=10"`
	Notes              *string `json:"notes,omitempty"`
}

type AcceptanceExerciseResponse struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	VideoWatched       bool      `json:"video_watched"`
	WatchedAt          string    `json:"watched_at,omitempty"`
	UnderstandingLevel *int      `json:"understanding_level,omitempty"`
	Notes              string    `json:"notes,omitempty"`
	DayNumber          *int      `json:"day_number,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type AcceptanceExerciseListRequest struct {
	PaginatedRequest
	FilterRequest
	VideoWatched *bool `query:"video_watched" form:"video_watched"`
}

type AcceptanceExerciseListResponse struct {
	PaginatedResponse[AcceptanceExerciseResponse]
}
