package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name   string `json:"category_name"`
	UserID uint
	Tasks  []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type CreateCategoryReq struct {
	Name   string `json:"category_name"`
	UserID uint
}

type UpdateCategoryReq struct {
	Name   string `json:"category_name"`
	UserID uint
}