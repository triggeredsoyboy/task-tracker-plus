package repository

import (
	"first-project/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *model.Category) error
	CategoryList(id int) ([]model.Category, error)
	UpdateCategory(id int, category *model.UpdateCategoryReq) error
	DeleteCategory(id int) error
	GetByID(id int) (*model.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) CreateCategory(category *model.Category) error {
	err := r.db.Create(&category).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) CategoryList(id int) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Preload("Tasks").Where("user_id = ?", id).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

func (r *categoryRepository) UpdateCategory(id int, category *model.UpdateCategoryReq) error {
	err := r.db.Table("categories").Where("id = ?", id).Updates(&category).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) DeleteCategory(id int) error {
	err := r.db.Where("id = ?", id).Delete(&model.Category{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) GetByID(id int) (*model.Category, error) {
	var category model.Category
	err := r.db.Preload("Tasks").Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}