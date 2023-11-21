package service

import (
	"first-project/model"
	repo "first-project/repository"
)

type CategoryService interface {
	CreateCategory(category *model.Category) error
	CategoryList(id int) ([]model.Category, error)
	UpdateCategory(id int, category *model.UpdateCategoryReq) error
	DeleteCategory(id int) error
	GetByID(id int) (*model.Category, error)
}

type categoryService struct {
	categoryRepository repo.CategoryRepository
}

func NewCategoryService(categoryRepository repo.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository}
}

func (s *categoryService) CreateCategory(category *model.Category) error {
	err := s.categoryRepository.CreateCategory(category)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) CategoryList(id int) ([]model.Category, error) {
	categories, err := s.categoryRepository.CategoryList(id)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *categoryService) UpdateCategory(id int, category *model.UpdateCategoryReq) error {
	err := s.categoryRepository.UpdateCategory(id, category)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) DeleteCategory(id int) error {
	err := s.categoryRepository.DeleteCategory(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) GetByID(id int) (*model.Category, error) {
	category, err := s.categoryRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}