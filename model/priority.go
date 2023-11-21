package model

import "gorm.io/gorm"

type Priority struct {
	gorm.Model
	PriorityClass string `json:"priority_class"`
	Tasks         []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}