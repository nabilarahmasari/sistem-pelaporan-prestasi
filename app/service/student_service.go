package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type StudentService struct {
	studentRepo  repository.StudentRepository
	lecturerRepo repository.LecturerRepository
	userRepo     repository.UserRepository
	validate     *validator.Validate
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	userRepo repository.UserRepository,
) *StudentService {
	return &StudentService{
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
		userRepo:     userRepo,
		validate:     validator.New(),
	}
}

//
// ==================== GET ALL STUDENTS (GET /students) ======================
//

func (s *StudentService) GetAllStudents(c *fiber.Ctx) error {
	// Implementasi pagination jika perlu (skip untuk simplicity)
	students, err := s.studentRepo.GetAll(100, 0)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch students",
		})
	}

	var responses []model.StudentResponse
	for _, student := range students {
		responses = append(responses, model.StudentResponse{
			ID:           student.ID,
			UserID:       student.UserID,
			StudentID:    student.StudentID,
			ProgramStudy: student.ProgramStudy,
			AcademicYear: student.AcademicYear,
			AdvisorID:    student.AdvisorID,
			CreatedAt:    student.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   responses,
	})
}

//
// ==================== GET STUDENT BY ID (GET /students/:id) ======================
//

func (s *StudentService) GetStudentByID(c *fiber.Ctx) error {
	studentID := c.Params("id")

	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	response := model.StudentResponse{
		ID:           student.ID,
		UserID:       student.UserID,
		StudentID:    student.StudentID,
		ProgramStudy: student.ProgramStudy,
		AcademicYear: student.AcademicYear,
		AdvisorID:    student.AdvisorID,
		CreatedAt:    student.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== SET ADVISOR (PUT /students/:id/advisor) ======================
//

func (s *StudentService) SetAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Cari student
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	// Parse request
	req := new(model.SetAdvisorRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Validasi advisor exists
	_, err = s.lecturerRepo.FindByID(req.AdvisorID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "advisor (lecturer) not found",
		})
	}

	// Update advisor
	if err := s.studentRepo.SetAdvisor(studentID, req.AdvisorID); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to set advisor",
		})
	}

	// Refresh student data
	student, _ = s.studentRepo.FindByID(studentID)

	response := model.StudentResponse{
		ID:           student.ID,
		UserID:       student.UserID,
		StudentID:    student.StudentID,
		ProgramStudy: student.ProgramStudy,
		AcademicYear: student.AcademicYear,
		AdvisorID:    student.AdvisorID,
		CreatedAt:    student.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "advisor set successfully",
		Data:    response,
	})
}
