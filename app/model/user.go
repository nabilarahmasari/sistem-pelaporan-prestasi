package model

import "time"

// ===================== USER ENTITY ========================
// Representasi tabel "users" di database

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` 
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ===================== USER CREATE DTO =====================
// Dipakai untuk endpoint: POST /users

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
}

// ===================== USER UPDATE DTO =====================
// Dipakai untuk endpoint: PUT /users/:id

type UserUpdateRequest struct {
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
}

// ===================== USER RESPONSE =======================
// Dipakai untuk mengembalikan informasi user ke frontend
// Tanpa password dan lebih ringan

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	Permissions []string `json:"permissions"` 
}
