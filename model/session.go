package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Token  string    `json:"token"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
}