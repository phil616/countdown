package repository

import (
	"ops-timer-backend/internal/model"

	"gorm.io/gorm"
)

type UnitLogRepository struct {
	db *gorm.DB
}

func NewUnitLogRepository(db *gorm.DB) *UnitLogRepository {
	return &UnitLogRepository{db: db}
}

func (r *UnitLogRepository) Create(log *model.UnitLog) error {
	return r.db.Create(log).Error
}

func (r *UnitLogRepository) ListByUnitID(unitID string, page, pageSize int) ([]model.UnitLog, int64, error) {
	var logs []model.UnitLog
	var total int64

	query := r.db.Model(&model.UnitLog{}).Where("unit_id = ?", unitID)
	query.Count(&total)

	err := query.Order("operated_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}
