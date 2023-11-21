package repository

import (
	"first-project/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	GetUserByID(id int) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(id int) (model.User, error) {
	var user model.User
	err := r.db.Preload("Categories.Tasks").Preload("Tasks").Where("id = ?", id).First(&user).Error
	if err != nil {
		return model.User{}, nil
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Preload("Categories.Tasks").Preload("Tasks").Where("email = ?", email).First(&user).Error
	if err != nil {
		return model.User{}, nil
	}

	return user, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.Omit(clause.Associations).Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}