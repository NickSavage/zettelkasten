package models

import (
	"time"
)

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    int       `json:"user_id"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EditTagParams struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}