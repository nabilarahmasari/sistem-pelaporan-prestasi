package service

import (
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type LecturerService struct {
	lecturerRepo repository.LecturerRepository
	studentRepo  repository.StudentRepository
	userRepo     repository.UserRepository
	validate     *validator.Validate
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
	userRepo repository.UserRepository,
) *LecturerService {
	return &LecturerService{
		lecturerRepo: lecturerRepo,
		studentRepo:  studentRepo,
		userRepo:     userRepo,
		validate:     validator.New(),
	}
}

//
// ==================== GET ALL LECTURERS (GET /lecturers) ======================
// SRS Section 5.5: GET /api/v1/lecturers
//

func (s *LecturerService) GetAllLecturers(c *fiber.Ctx) error {
	// Parse query params for pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get lecturers
	lecturers, err := s.lecturerRepo.GetAll(pageSize, offset)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch lecturers",
		})
	}

	// Count total
	total, err := s.lecturerRepo.CountAll()
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count lecturers",
		})
	}

	// Build responses with user info
	var responses []map[string]interface{}
	for _, lecturer := range lecturers {
		// Get user info
		user, _ := s.userRepo.FindByID(lecturer.UserID)

		lecturerData := map[string]interface{}{
			"id":          lecturer.ID,
			"user_id":     lecturer.UserID,
			"lecturer_id": lecturer.LecturerID,
			"department":  lecturer.Department,
			"created_at":  lecturer.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if user != nil {
			lecturerData["username"] = user.Username
			lecturerData["full_name"] = user.FullName
			lecturerData["email"] = user.Email
		}

		responses = append(responses, lecturerData)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"lecturers":   responses,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

//
// ==================== GET LECTURER ADVISEES (GET /lecturers/:id/advisees) ======================
// SRS Section 5.5: GET /api/v1/lecturers/:id/advisees
// FR-006: Dosen wali melihat daftar mahasiswa bimbingannya
//

func (s *LecturerService) GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	// Get user dari context untuk authorization
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Check if lecturer exists
	lecturer, err := s.lecturerRepo.FindByID(lecturerID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "lecturer not found",
		})
	}

	// Authorization: Dosen hanya bisa lihat advisees sendiri, Admin bisa lihat semua
	if claims.Role == "Dosen Wali" {
		currentLecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		if currentLecturer == nil || currentLecturer.ID != lecturerID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: can only view your own advisees",
			})
		}
	}
	// Admin dapat akses semua

	// Get all students
	allStudents, err := s.studentRepo.GetAll(1000, 0) // Get all for filtering
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch students",
		})
	}

	// Filter students by advisor_id
	var advisees []map[string]interface{}
	for _, student := range allStudents {
		if student.AdvisorID != nil && *student.AdvisorID == lecturerID {
			// Get user info
			user, _ := s.userRepo.FindByID(student.UserID)

			adviseeData := map[string]interface{}{
				"id":            student.ID,
				"student_id":    student.StudentID,
				"program_study": student.ProgramStudy,
				"academic_year": student.AcademicYear,
				"created_at":    student.CreatedAt.Format("2006-01-02 15:04:05"),
			}

			if user != nil {
				adviseeData["full_name"] = user.FullName
				adviseeData["email"] = user.Email
			}

			advisees = append(advisees, adviseeData)
		}
	}

	// Get lecturer user info
	lecturerUser, _ := s.userRepo.FindByID(lecturer.UserID)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"lecturer": map[string]interface{}{
				"id":          lecturer.ID,
				"lecturer_id": lecturer.LecturerID,
				"full_name":   lecturerUser.FullName,
				"department":  lecturer.Department,
			},
			"advisees": advisees,
			"total":    len(advisees),
		},
	})
}
