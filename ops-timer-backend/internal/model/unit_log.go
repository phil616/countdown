package model

import "time"

type UnitLog struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	UnitID      string    `gorm:"size:36;not null;index" json:"unit_id"`
	Delta       float64   `gorm:"not null" json:"delta"`
	ValueBefore float64   `gorm:"not null" json:"value_before"`
	ValueAfter  float64   `gorm:"not null" json:"value_after"`
	Note        string    `gorm:"size:512" json:"note"`
	OperatedAt  time.Time `gorm:"not null" json:"operated_at"`
}
