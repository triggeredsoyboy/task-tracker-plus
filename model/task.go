package model

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title		string    `json:"title"`
	Description string    `json:"description"`
	Deadline	time.Time `json:"deadline"`
	Status		string    `json:"status"`
	UserID		uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	PriorityID  uint      `json:"priority_id"`
}

type CreateTaskReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	PriorityID  uint      `json:"priority_id"`
}

type UpdateTaskReq struct {
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	PriorityID  uint      `json:"priority_id"`
}

type TaskByCategory struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	Deadline    time.Time `json:"deadline"`
	Status      string    `json:"status"`
	UserID      uint      `json:"user_id"`
	PriorityID  uint      `json:"priority_id"`
}