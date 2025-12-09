package model

import "time"

// ===================== STUDENT ENTITY ========================
// Representasi tabel "students" di database

type Student struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	StudentID    string    `json:"student_id" db:"student_id"`
	ProgramStudy string    `json:"program_study" db:"program_study"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	AdvisorID    *string   `json:"advisor_id" db:"advisor_id"` // nullable
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ===================== STUDENT PROFILE DTO ====================
// Dipakai saat create/update user dengan role mahasiswa

type StudentProfileRequest struct {
	StudentID    string  `json:"student_id" validate:"required"`
	ProgramStudy string  `json:"program_study" validate:"required"`
	AcademicYear string  `json:"academic_year" validate:"required"`
	AdvisorID    *string `json:"advisor_id,omitempty"` // optional saat create
}

// ===================== STUDENT RESPONSE =======================
// Response untuk menampilkan data student

type StudentResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	StudentID    string  `json:"student_id"`
	ProgramStudy string  `json:"program_study"`
	AcademicYear string  `json:"academic_year"`
	AdvisorID    *string `json:"advisor_id,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

// ===================== SET ADVISOR REQUEST ====================
// Untuk endpoint PUT /students/:id/advisor

type SetAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" validate:"required"`
}
