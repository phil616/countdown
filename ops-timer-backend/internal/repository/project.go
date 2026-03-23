package repository

import (
	"ops-timer-backend/internal/model"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) FindByID(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) List(status string, page, pageSize int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	err := query.Order("sort_order ASC, created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&projects).Error

	return projects, total, err
}

func (r *ProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id string) error {
	return r.db.Delete(&model.Project{}, "id = ?", id).Error
}

func (r *ProjectRepository) ClearProjectUnits(projectID string) error {
	return r.db.Model(&model.Unit{}).Where("project_id = ?", projectID).
		Update("project_id", nil).Error
}
