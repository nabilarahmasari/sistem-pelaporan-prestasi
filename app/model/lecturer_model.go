package model

import "time"

// ===================== LECTURER ENTITY ========================
// Representasi tabel "lecturers" di database

type Lecturer struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	LecturerID string    `json:"lecturer_id" db:"lecturer_id"`
	Department string    `json:"department" db:"department"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ===================== LECTURER PROFILE DTO ===================
// Dipakai saat create/update user dengan role dosen

type LecturerProfileRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
	Department string `json:"department" validate:"required"`
}

// ===================== LECTURER RESPONSE ======================
// Response untuk menampilkan data lecturer

type LecturerResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
	CreatedAt  string `json:"created_at"`
}
