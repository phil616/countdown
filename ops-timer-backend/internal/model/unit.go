package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type JSONFloatArray []float64

func (a JSONFloatArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *JSONFloatArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return nil
	}
	return json.Unmarshal([]byte(s), a)
}

type JSONIntArray []int

func (a JSONIntArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *JSONIntArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return nil
	}
	return json.Unmarshal([]byte(s), a)
}

type JSONStringArray []string

func (a JSONStringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return nil
	}
	return json.Unmarshal([]byte(s), a)
}

const (
	UnitTypeTimeCountdown  = "time_countdown"
	UnitTypeTimeCountup    = "time_countup"
	UnitTypeCountCountdown = "count_countdown"
	UnitTypeCountCountup   = "count_countup"

	UnitStatusActive    = "active"
	UnitStatusPaused    = "paused"
	UnitStatusCompleted = "completed"
	UnitStatusArchived  = "archived"

	PriorityLow      = "low"
	PriorityNormal   = "normal"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

type Unit struct {
	ID          string         `gorm:"primaryKey;size:36" json:"id"`
	ProjectID   *string        `gorm:"size:36;index" json:"project_id"`
	Title       string         `gorm:"size:128;not null" json:"title"`
	Description string         `gorm:"size:2048" json:"description"`
	Type        string         `gorm:"size:20;not null" json:"type"`
	Status      string         `gorm:"size:20;not null;default:active" json:"status"`
	Priority    string         `gorm:"size:10;not null;default:normal" json:"priority"`
	Tags        JSONStringArray `gorm:"type:text" json:"tags"`
	Color       string         `gorm:"size:20" json:"color"`

	// Time-based fields
	TargetTime       *time.Time   `json:"target_time"`
	StartTime        *time.Time   `json:"start_time"`
	DisplayUnit      string       `gorm:"size:10;default:days" json:"display_unit"`
	RemindBeforeDays JSONIntArray `gorm:"type:text" json:"remind_before_days"`
	RemindAfterDays  JSONIntArray `gorm:"type:text" json:"remind_after_days"`

	// Count-based fields
	CurrentValue   *float64       `json:"current_value"`
	TargetValue    *float64       `json:"target_value"`
	Step           float64        `gorm:"default:1" json:"step"`
	UnitLabel      string         `gorm:"size:20" json:"unit_label"`
	AllowExceed    bool           `gorm:"default:false" json:"allow_exceed"`
	RemindOnValues JSONFloatArray `gorm:"type:text" json:"remind_on_values"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (u *Unit) IsTimeType() bool {
	return u.Type == UnitTypeTimeCountdown || u.Type == UnitTypeTimeCountup
}

func (u *Unit) IsCountType() bool {
	return u.Type == UnitTypeCountCountdown || u.Type == UnitTypeCountCountup
}

func (u *Unit) IsCountdown() bool {
	return u.Type == UnitTypeTimeCountdown || u.Type == UnitTypeCountCountdown
}
