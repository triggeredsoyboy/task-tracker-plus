package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fullname   string     `json:"fullname"`
	Email      string     `gorm:"unique" json:"email"`
	Password   string     `json:"-"`
	Categories []Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tasks      []Task     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type RegisterData struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}