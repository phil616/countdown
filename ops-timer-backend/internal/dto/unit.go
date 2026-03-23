package dto

import "time"

type CreateUnitRequest struct {
	ProjectID   *string    `json:"project_id"`
	Title       string     `json:"title" binding:"required,max=128"`
	Description string     `json:"description" binding:"max=2048"`
	Type        string     `json:"type" binding:"required,oneof=time_countdown time_countup count_countdown count_countup"`
	Status      string     `json:"status" binding:"omitempty,oneof=active paused completed archived"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low normal high critical"`
	Tags        []string   `json:"tags"`
	Color       string     `json:"color"`

	// Time fields
	TargetTime       *time.Time `json:"target_time"`
	StartTime        *time.Time `json:"start_time"`
	DisplayUnit      string     `json:"display_unit" binding:"omitempty,oneof=days hours minutes seconds"`
	RemindBeforeDays []int      `json:"remind_before_days"`
	RemindAfterDays  []int      `json:"remind_after_days"`

	// Count fields
	CurrentValue   *float64  `json:"current_value"`
	TargetValue    *float64  `json:"target_value"`
	Step           *float64  `json:"step"`
	UnitLabel      string    `json:"unit_label"`
	AllowExceed    *bool     `json:"allow_exceed"`
	RemindOnValues []float64 `json:"remind_on_values"`
}

type UpdateUnitRequest struct {
	ProjectID   *string `json:"project_id"`
	ClearProject bool   `json:"clear_project"` // 设为 true 时将单元从项目中移除
	Title       string  `json:"title" binding:"omitempty,max=128"`
	Description *string    `json:"description"`
	Status      string     `json:"status" binding:"omitempty,oneof=active paused completed archived"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low normal high critical"`
	Tags        []string   `json:"tags"`
	Color       *string    `json:"color"`

	TargetTime       *time.Time `json:"target_time"`
	StartTime        *time.Time `json:"start_time"`
	DisplayUnit      string     `json:"display_unit" binding:"omitempty,oneof=days hours minutes seconds"`
	RemindBeforeDays []int      `json:"remind_before_days"`
	RemindAfterDays  []int      `json:"remind_after_days"`

	CurrentValue   *float64  `json:"current_value"`
	TargetValue    *float64  `json:"target_value"`
	Step           *float64  `json:"step"`
	UnitLabel      *string   `json:"unit_label"`
	AllowExceed    *bool     `json:"allow_exceed"`
	RemindOnValues []float64 `json:"remind_on_values"`
}

type UpdateUnitStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active paused completed archived"`
}

type StepRequest struct {
	Direction string `json:"direction" binding:"required,oneof=up down"`
	Note      string `json:"note"`
}

type SetValueRequest struct {
	Value float64 `json:"value" binding:"required"`
	Note  string  `json:"note"`
}

type UnitResponse struct {
	ID          string    `json:"id"`
	ProjectID   *string   `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Tags        []string  `json:"tags"`
	Color       string    `json:"color"`

	TargetTime       *time.Time `json:"target_time,omitempty"`
	StartTime        *time.Time `json:"start_time,omitempty"`
	DisplayUnit      string     `json:"display_unit,omitempty"`
	RemindBeforeDays []int      `json:"remind_before_days,omitempty"`
	RemindAfterDays  []int      `json:"remind_after_days,omitempty"`

	CurrentValue   *float64  `json:"current_value,omitempty"`
	TargetValue    *float64  `json:"target_value,omitempty"`
	Step           float64   `json:"step,omitempty"`
	UnitLabel      string    `json:"unit_label,omitempty"`
	AllowExceed    bool      `json:"allow_exceed,omitempty"`
	RemindOnValues []float64 `json:"remind_on_values,omitempty"`

	// Computed fields
	RemainingSeconds *float64 `json:"remaining_seconds,omitempty"`
	ElapsedSeconds   *float64 `json:"elapsed_seconds,omitempty"`
	Progress         *float64 `json:"progress,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UnitSummary struct {
	TotalActive    int64 `json:"total_active"`
	TotalPaused    int64 `json:"total_paused"`
	TotalCompleted int64 `json:"total_completed"`
	TotalArchived  int64 `json:"total_archived"`
	ExpiringCount  int64 `json:"expiring_count"`
	ExpiredCount   int64 `json:"expired_count"`
}

type UnitLogResponse struct {
	ID          string    `json:"id"`
	UnitID      string    `json:"unit_id"`
	Delta       float64   `json:"delta"`
	ValueBefore float64   `json:"value_before"`
	ValueAfter  float64   `json:"value_after"`
	Note        string    `json:"note"`
	OperatedAt  time.Time `json:"operated_at"`
}

type UnitQueryParams struct {
	Type      string `form:"type"`
	Status    string `form:"status"`
	ProjectID string `form:"project_id"`
	Tags      string `form:"tags"`
	Priority  string `form:"priority"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
	Q         string `form:"q"`
	Page      int    `form:"page,default=1"`
	PageSize  int    `form:"page_size,default=20"`
}
