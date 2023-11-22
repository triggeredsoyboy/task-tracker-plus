package repository

import (
	"first-project/model"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *model.Task) error
	TaskList(id int) ([]model.Task, error)
	UpdateTask(id int, task *model.UpdateTaskReq) error
	DeleteTask(id int) error
	GetByID(id int) (*model.Task, error)
	GetByCategory(id int) ([]model.TaskByCategory, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *taskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) CreateTask(task *model.Task) error {
	err := r.db.Create(&task).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepository) TaskList(id int) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("user_id = ?", id).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) UpdateTask(id int, task *model.UpdateTaskReq) error {
	err := r.db.Table("tasks").Where("id = ?", id).Updates(&task).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepository) DeleteTask(id int) error {
	err := r.db.Where("id = ?", id).Delete(&model.Task{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepository) GetByID(id int) (*model.Task, error) {
	var task model.Task
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) GetByCategory(id int) ([]model.TaskByCategory, error) {
	var tasks []model.TaskByCategory
	err := r.db.Table("tasks").Select("tasks.id, tasks.title, tasks.description, categories.name AS category, tasks.created_at, tasks.deadline, tasks.status, tasks.priority_id").Joins("LEFT JOIN categories ON tasks.category_id = categories.id").Where("tasks.deleted_at IS NULL").Where("tasks.user_id = ?", id).Scan(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}