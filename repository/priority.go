package repository

import (
	"first-project/model"

	"gorm.io/gorm"
)

type PriorityRepository interface {
	GetByID(id int) (*model.Priority, error)
}

type priorityRepository struct {
	db *gorm.DB
}

func NewPriorityRepo(db *gorm.DB) *priorityRepository {
	return &priorityRepository{db}
}

func (r *priorityRepository) GetByID(id int) (*model.Priority, error) {
	var priority model.Priority
	err := r.db.Where("id = ?", id).First(&priority).Error
	if err != nil {
		return nil, err
	}

	return &priority, nil
}