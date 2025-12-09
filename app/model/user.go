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

// ===================== USER CREATE REQUEST =====================
// Dipakai untuk endpoint: POST /api/v1/users (Admin only)
// Support create dengan profile student/lecturer sekaligus

type UserCreateRequest struct {
	Username       string                  `json:"username" validate:"required"`
	Email          string                  `json:"email" validate:"required,email"`
	Password       string                  `json:"password" validate:"required,min=8"`
	FullName       string                  `json:"full_name" validate:"required"`
	RoleName       string                  `json:"role_name" validate:"required"` // "Admin", "Mahasiswa", "Dosen Wali"
	StudentProfile *StudentProfileRequest  `json:"student_profile,omitempty"`     // jika role = Mahasiswa
	LecturerProfile *LecturerProfileRequest `json:"lecturer_profile,omitempty"`   // jika role = Dosen Wali
}

// ===================== USER UPDATE DTO =====================
// Dipakai untuk endpoint: PUT /users/:id

type UserUpdateRequest struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	FullName string `json:"full_name,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"` // untuk activate/deactivate user
}

// ===================== ASSIGN ROLE REQUEST ===================
// Untuk endpoint PUT /users/:id/role

type AssignRoleRequest struct {
	RoleName string `json:"role_name" validate:"required"`
}

// ===================== USER RESPONSE =======================
// Dipakai untuk mengembalikan informasi user ke frontend
// Tanpa password dan lebih ringan

type UserResponse struct {
	ID              string            `json:"id"`
	Username        string            `json:"username"`
	Email           string            `json:"email"`
	FullName        string            `json:"full_name"`
	Role            string            `json:"role"`
	IsActive        bool              `json:"is_active"`
	CreatedAt       string            `json:"created_at"`
	Permissions     []string          `json:"permissions,omitempty"` 
	StudentProfile  *StudentResponse  `json:"student_profile,omitempty"`  // jika role = Mahasiswa
	LecturerProfile *LecturerResponse `json:"lecturer_profile,omitempty"` // jika role = Dosen Wali
}

// ===================== USER LIST RESPONSE ====================
// Untuk endpoint GET /users dengan pagination

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}