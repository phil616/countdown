package repository

import (
	"ops-timer-backend/internal/model"

	"gorm.io/gorm"
)

type TodoGroupRepository struct {
	db *gorm.DB
}

func NewTodoGroupRepository(db *gorm.DB) *TodoGroupRepository {
	return &TodoGroupRepository{db: db}
}

func (r *TodoGroupRepository) Create(group *model.TodoGroup) error {
	return r.db.Create(group).Error
}

func (r *TodoGroupRepository) FindByID(id string) (*model.TodoGroup, error) {
	var group model.TodoGroup
	err := r.db.First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *TodoGroupRepository) List() ([]model.TodoGroup, error) {
	var groups []model.TodoGroup
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&groups).Error
	if err != nil {
		return nil, err
	}

	for i := range groups {
		var count int64
		r.db.Model(&model.Todo{}).Where("group_id = ?", groups[i].ID).Count(&count)
		groups[i].TodoCount = count
	}

	return groups, nil
}

func (r *TodoGroupRepository) Update(group *model.TodoGroup) error {
	return r.db.Save(group).Error
}

func (r *TodoGroupRepository) Delete(id string) error {
	return r.db.Delete(&model.TodoGroup{}, "id = ?", id).Error
}
