package models

import (
	"time"
)

type Table struct {
	ID         string    `json:"id" db:"id"`
	LocationID string    `json:"location_id" db:"location_id"`
	Code       string    `json:"code" db:"code"`
	Seats      int       `json:"seats" db:"seats"`
	Status     string    `json:"status" db:"status"` // free, occupied, reserved
	IsActive   bool      `json:"is_active" db:"is_active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTableRequest struct {
	LocationID string `json:"location_id" binding:"required"`
	Code       string `json:"code" binding:"required"`
	Seats      int    `json:"seats" binding:"required,min=1,max=20"`
}

type UpdateTableRequest struct {
	Code     string `json:"code"`
	Seats    int    `json:"seats" binding:"min=1,max=20"`
	Status   string `json:"status" binding:"omitempty,oneof=free occupied reserved"`
	IsActive *bool  `json:"is_active"`
}
