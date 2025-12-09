package repository

import (
	"database/sql"
	"fmt"
	"project_uas/app/model"
	"time"
)

type StudentRepository interface {
	Create(student *model.Student) error
	FindByUserID(userID string) (*model.Student, error)
	FindByID(id string) (*model.Student, error)
	FindByStudentID(studentID string) (*model.Student, error)
	Update(student *model.Student) error
	Delete(id string) error
	SetAdvisor(studentID string, advisorID string) error
	GetAll(limit, offset int) ([]model.Student, error)
	CountAll() (int, error)
}

type studentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db}
}

// Create - Insert student baru
// SESUAI MIGRATION: id = user_id, academic_year = INT
func (r *studentRepository) Create(student *model.Student) error {
	// ID = user_id (sesuai migration: id REFERENCES users(id))
	student.ID = student.UserID
	student.CreatedAt = time.Now()

	// Convert academic_year string to int
	var academicYear int
	if student.AcademicYear != "" {
		fmt.Sscanf(student.AcademicYear, "%d", &academicYear)
	}

	query := `
		INSERT INTO students (id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query,
		student.ID,           // id = user_id
		student.StudentID,    // student_id (NIM)
		student.ProgramStudy, // program_study
		academicYear,         // academic_year (INT)
		student.AdvisorID,    // advisor_id
		student.CreatedAt,
	)
	return err
}

// FindByUserID - Cari student berdasarkan user_id (id = user_id)
func (r *studentRepository) FindByUserID(userID string) (*model.Student, error) {
	student := &model.Student{}
	var academicYear int
	
	query := `
		SELECT id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`
	
	err := r.db.QueryRow(query, userID).Scan(
		&student.ID,
		&student.StudentID,
		&student.ProgramStudy,
		&academicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	student.UserID = student.ID // id = user_id
	student.AcademicYear = fmt.Sprintf("%d", academicYear)
	return student, nil
}

// FindByID - Cari student berdasarkan ID
func (r *studentRepository) FindByID(id string) (*model.Student, error) {
	student := &model.Student{}
	var academicYear int
	
	query := `
		SELECT id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.StudentID,
		&student.ProgramStudy,
		&academicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	student.UserID = student.ID
	student.AcademicYear = fmt.Sprintf("%d", academicYear)
	return student, nil
}

// FindByStudentID - Cari student berdasarkan student_id (NIM)
func (r *studentRepository) FindByStudentID(studentID string) (*model.Student, error) {
	student := &model.Student{}
	var academicYear int
	
	query := `
		SELECT id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		WHERE student_id = $1
	`
	
	err := r.db.QueryRow(query, studentID).Scan(
		&student.ID,
		&student.StudentID,
		&student.ProgramStudy,
		&academicYear,
		&student.AdvisorID,
		&student.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	student.UserID = student.ID
	student.AcademicYear = fmt.Sprintf("%d", academicYear)
	return student, nil
}

// Update - Update data student
func (r *studentRepository) Update(student *model.Student) error {
	var academicYear int
	fmt.Sscanf(student.AcademicYear, "%d", &academicYear)
	
	query := `
		UPDATE students
		SET program_study = $1, academic_year = $2, advisor_id = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query,
		student.ProgramStudy,
		academicYear,
		student.AdvisorID,
		student.ID,
	)
	return err
}

// Delete - Hapus student
func (r *studentRepository) Delete(id string) error {
	query := `DELETE FROM students WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// SetAdvisor - Set dosen pembimbing untuk student
func (r *studentRepository) SetAdvisor(studentID string, advisorID string) error {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, advisorID, studentID)
	return err
}

// GetAll - Ambil semua students dengan pagination
func (r *studentRepository) GetAll(limit, offset int) ([]model.Student, error) {
	query := `
		SELECT id, student_id, program_study, academic_year, advisor_id, created_at
		FROM students
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		var academicYear int
		
		err := rows.Scan(
			&s.ID,
			&s.StudentID,
			&s.ProgramStudy,
			&academicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		)
		if err != nil {
			continue
		}
		
		s.UserID = s.ID
		s.AcademicYear = fmt.Sprintf("%d", academicYear)
		students = append(students, s)
	}
	return students, nil
}

// CountAll - Hitung total students
func (r *studentRepository) CountAll() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM students`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}