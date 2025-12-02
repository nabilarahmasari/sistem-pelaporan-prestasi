package model

import "time"

type Role struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"` // ⭐ TAMBAHKAN INI
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"` // ⭐ TAMBAHKAN INI JUGA (opsional)
}