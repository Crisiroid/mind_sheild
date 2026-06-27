package schemas

import "time"

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PaginatedRequest struct {
	Page      int    `query:"page" form:"page" default:"1"`
	PageSize  int    `query:"page_size" form:"page_size" default:"20"`
	SortBy    string `query:"sort_by" form:"sort_by" default:"created_at"`
	SortOrder string `query:"sort_order" form:"sort_order" default:"desc"`
}

type PaginatedResponse[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Pages    int   `json:"pages"`
}

type FilterRequest struct {
	DateFrom *time.Time `query:"date_from" form:"date_from"`
	DateTo   *time.Time `query:"date_to" form:"date_to"`
	UserID   string     `query:"user_id" form:"user_id"`
	Search   string     `query:"search" form:"search"`
}

type DateRangeFilter struct {
	From time.Time `json:"from" validate:"required"`
	To   time.Time `json:"to" validate:"required"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

type ReportMeta struct {
	GeneratedAt  time.Time `json:"generated_at"`
	DateFrom     time.Time `json:"date_from"`
	DateTo       time.Time `json:"date_to"`
	TotalRecords int64     `json:"total_records"`
}

type TimeSeriesData struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label,omitempty"`
}

type DistributionItem struct {
	Label string  `json:"label"`
	Count int64   `json:"count"`
	Value float64 `json:"value,omitempty"`
}

type RankingItem struct {
	Rank  int     `json:"rank"`
	Label string  `json:"label"`
	Count int64   `json:"count"`
	Value float64 `json:"value,omitempty"`
}
