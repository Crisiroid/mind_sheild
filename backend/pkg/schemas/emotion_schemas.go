package schemas

import "time"

type EmotionInteractionCreateRequest struct {
	UserID                string   `json:"user_id" validate:"required"`
	SideClicked           string   `json:"side_clicked" validate:"required,oneof=thought body behavior"`
	ThoughtAccountsViewed []string `json:"thought_accounts_viewed,omitempty"`
	VibrationTriggered    bool     `json:"vibration_triggered,omitempty"`
	DayNumber             *int     `json:"day_number,omitempty"`
}

type EmotionInteractionResponse struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	InteractionDate       time.Time `json:"interaction_date"`
	SideClicked           string    `json:"side_clicked"`
	ThoughtAccountsViewed []string  `json:"thought_accounts_viewed,omitempty"`
	VibrationTriggered    bool      `json:"vibration_triggered"`
	DayNumber             *int      `json:"day_number,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type EmotionInteractionListRequest struct {
	PaginatedRequest
	FilterRequest
	SideClicked *string `query:"side_clicked" form:"side_clicked"`
}

type EmotionInteractionUpdateRequest struct {
	SideClicked           *string  `json:"side_clicked,omitempty" validate:"omitempty,oneof=thought body behavior"`
	ThoughtAccountsViewed []string `json:"thought_accounts_viewed,omitempty"`
	VibrationTriggered    *bool    `json:"vibration_triggered,omitempty"`
}

type EmotionInteractionListResponse struct {
	PaginatedResponse[EmotionInteractionResponse]
}

type StressEventCreateRequest struct {
	UserID               string `json:"user_id" validate:"required"`
	SituationType        string `json:"situation_type" validate:"required"`
	SituationDescription string `json:"situation_description,omitempty"`
	IntensityLevel       int    `json:"intensity_level" validate:"required,min=1,max=10"`
	Location             string `json:"location,omitempty"`
	DayNumber            *int   `json:"day_number,omitempty"`
}

type StressEventUpdateRequest struct {
	IntensityLevel       *int    `json:"intensity_level,omitempty" validate:"omitempty,min=1,max=10"`
	SituationDescription *string `json:"situation_description,omitempty"`
}

type StressEventResponse struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	EventTimestamp       time.Time `json:"event_timestamp"`
	SituationType        string    `json:"situation_type"`
	SituationDescription string    `json:"situation_description,omitempty"`
	IntensityLevel       int       `json:"intensity_level"`
	Location             string    `json:"location,omitempty"`
	DayNumber            *int      `json:"day_number,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type StressEventListRequest struct {
	PaginatedRequest
	FilterRequest
	SituationType *string `query:"situation_type" form:"situation_type"`
	IntensityMin  *int    `query:"intensity_min" form:"intensity_min"`
	IntensityMax  *int    `query:"intensity_max" form:"intensity_max"`
}

type StressEventListResponse struct {
	PaginatedResponse[StressEventResponse]
}

type BodyTensionCreateRequest struct {
	UserID           string `json:"user_id" validate:"required"`
	BodyRegions      string `json:"body_regions" validate:"required"`
	OverallIntensity *int   `json:"overall_intensity,omitempty" validate:"omitempty,min=1,max=10"`
	SeverityColor    string `json:"severity_color,omitempty"`
	Notes            string `json:"notes,omitempty"`
	DayNumber        *int   `json:"day_number,omitempty"`
}

type BodyTensionUpdateRequest struct {
	BodyRegions      *string `json:"body_regions,omitempty"`
	OverallIntensity *int    `json:"overall_intensity,omitempty" validate:"omitempty,min=1,max=10"`
	Notes            *string `json:"notes,omitempty"`
}

type BodyTensionResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	MappingDate      time.Time `json:"mapping_date"`
	BodyRegions      string    `json:"body_regions"`
	OverallIntensity *int      `json:"overall_intensity,omitempty"`
	SeverityColor    string    `json:"severity_color,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	DayNumber        *int      `json:"day_number,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type BodyTensionListRequest struct {
	PaginatedRequest
	FilterRequest
	IntensityMin *int `query:"intensity_min" form:"intensity_min"`
	IntensityMax *int `query:"intensity_max" form:"intensity_max"`
}

type BodyTensionListResponse struct {
	PaginatedResponse[BodyTensionResponse]
}
