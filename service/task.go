package service

import (
	"first-project/model"
	repo "first-project/repository"
)

type TaskService interface {
	CreateTask(task *model.Task) error
	TaskList(id int) ([]model.Task, error)
	UpdateTask(id int, task *model.UpdateTaskReq) error
	DeleteTask(id int) error
	GetByID(id int) (*model.Task, error)
	GetByCategory(id int) ([]model.TaskByCategory, error)
}

type taskService struct {
	taskRepository repo.TaskRepository
}

func NewTaskService(taskRepository repo.TaskRepository) TaskService {
	return &taskService{taskRepository}
}

func (s *taskService) CreateTask(task *model.Task) error {
	err := s.taskRepository.CreateTask(task)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) TaskList(id int) ([]model.Task, error) {
	tasks, err := s.taskRepository.TaskList(id)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *taskService) UpdateTask(id int, task *model.UpdateTaskReq) error {
	err := s.taskRepository.UpdateTask(id, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) DeleteTask(id int) error {
	err := s.taskRepository.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) GetByID(id int) (*model.Task, error) {
	task, err := s.taskRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetByCategory(id int) ([]model.TaskByCategory, error) {
	tasks, err := s.taskRepository.GetByCategory(id)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}