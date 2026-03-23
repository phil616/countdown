package model

import "time"

const (
	NotificationLevelInfo     = "info"
	NotificationLevelWarning  = "warning"
	NotificationLevelCritical = "critical"
)

type Notification struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	UnitID      string     `gorm:"size:36;index" json:"unit_id"`
	Level       string     `gorm:"size:10;not null" json:"level"`
	Message     string     `gorm:"size:512;not null" json:"message"`
	IsRead      bool       `gorm:"default:false" json:"is_read"`
	TriggeredAt time.Time  `gorm:"not null" json:"triggered_at"`
	ReadAt      *time.Time `json:"read_at"`

	Unit *Unit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}
