package service

import (
	"first-project/model"
	repo "first-project/repository"
)

type PriorityService interface {
	GetByID(id int) (*model.Priority, error)
}

type priorityService struct {
	priorityRepository repo.PriorityRepository
}

func NewPriorityService(priorityRepository repo.PriorityRepository) PriorityService {
	return &priorityService{priorityRepository}
}

func (s *priorityService) GetByID(id int) (*model.Priority, error) {
	priority, err := s.priorityRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return priority, nil
}