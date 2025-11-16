package models

import (
	"time"
)

type Location struct {
	ID        string    `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateLocationRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type UpdateLocationRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	IsActive *bool  `json:"is_active"`
}
