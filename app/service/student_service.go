package service

import (
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"project_uas/app/model"
	"project_uas/app/repository"
)

type StudentService struct {
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
	userRepo        repository.UserRepository
	achievementRepo repository.AchievementRepository
	validate        *validator.Validate
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	userRepo repository.UserRepository,
	achievementRepo repository.AchievementRepository,
) *StudentService {
	return &StudentService{
		studentRepo:     studentRepo,
		lecturerRepo:    lecturerRepo,
		userRepo:        userRepo,
		achievementRepo: achievementRepo,
		validate:        validator.New(),
	}
}

//
// ==================== GET ALL STUDENTS (GET /students) ======================
// SRS Section 5.5: GET /api/v1/students
//

func (s *StudentService) GetAllStudents(c *fiber.Ctx) error {
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

	// Get students with pagination
	students, err := s.studentRepo.GetAll(pageSize, offset)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch students",
		})
	}

	// Count total
	total, err := s.studentRepo.CountAll()
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count students",
		})
	}

	// Build responses with user info
	var responses []map[string]interface{}
	for _, student := range students {
		// Get user info
		user, _ := s.userRepo.FindByID(student.UserID)
		
		studentData := map[string]interface{}{
			"id":            student.ID,
			"user_id":       student.UserID,
			"student_id":    student.StudentID,
			"program_study": student.ProgramStudy,
			"academic_year": student.AcademicYear,
			"advisor_id":    student.AdvisorID,
			"created_at":    student.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if user != nil {
			studentData["username"] = user.Username
			studentData["full_name"] = user.FullName
			studentData["email"] = user.Email
		}

		// Get advisor info if exists
		if student.AdvisorID != nil {
			advisor, _ := s.lecturerRepo.FindByID(*student.AdvisorID)
			if advisor != nil {
				advisorUser, _ := s.userRepo.FindByID(advisor.UserID)
				if advisorUser != nil {
					studentData["advisor_name"] = advisorUser.FullName
				}
			}
		}

		responses = append(responses, studentData)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"students":    responses,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": totalPages,
		},
	})
}

//
// ==================== GET STUDENT BY ID (GET /students/:id) ======================
// SRS Section 5.5: GET /api/v1/students/:id
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

	// Get user info
	user, _ := s.userRepo.FindByID(student.UserID)

	response := map[string]interface{}{
		"id":            student.ID,
		"user_id":       student.UserID,
		"student_id":    student.StudentID,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
		"advisor_id":    student.AdvisorID,
		"created_at":    student.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user != nil {
		response["username"] = user.Username
		response["full_name"] = user.FullName
		response["email"] = user.Email
	}

	// Get advisor info if exists
	if student.AdvisorID != nil {
		advisor, _ := s.lecturerRepo.FindByID(*student.AdvisorID)
		if advisor != nil {
			advisorUser, _ := s.userRepo.FindByID(advisor.UserID)
			if advisorUser != nil {
				response["advisor"] = map[string]interface{}{
					"id":          advisor.ID,
					"lecturer_id": advisor.LecturerID,
					"full_name":   advisorUser.FullName,
					"department":  advisor.Department,
				}
			}
		}
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== GET STUDENT ACHIEVEMENTS (GET /students/:id/achievements) ======================
// SRS Section 5.5: GET /api/v1/students/:id/achievements
// Authorization: Admin dapat akses semua, Dosen Wali hanya mahasiswa bimbingannya, Mahasiswa hanya milik sendiri
//

func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Get user dari context
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Check if student exists
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "student not found",
		})
	}

	// Authorization check
	if claims.Role == "Mahasiswa" {
		// Mahasiswa hanya bisa akses prestasi sendiri
		currentStudent, _ := s.studentRepo.FindByUserID(claims.UserID)
		if currentStudent == nil || currentStudent.ID != studentID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden",
			})
		}
	} else if claims.Role == "Dosen Wali" {
		// Dosen hanya bisa akses prestasi mahasiswa bimbingannya
		lecturer, _ := s.lecturerRepo.FindByUserID(claims.UserID)
		if lecturer == nil || student.AdvisorID == nil || *student.AdvisorID != lecturer.ID {
			return c.Status(403).JSON(model.APIResponse{
				Status: "error",
				Error:  "forbidden: not your advisee",
			})
		}
	}
	// Admin dapat akses semua

	// Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	status := c.Query("status", "")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get achievements
	references, err := s.achievementRepo.GetReferencesByStudentID(studentID, status, pageSize, offset)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch achievements",
		})
	}

	total, err := s.achievementRepo.CountReferencesByStudentID(studentID, status)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count achievements",
		})
	}

	// Fetch details dari MongoDB
	var achievements []model.AchievementResponse
	for _, ref := range references {
		achievement, err := s.achievementRepo.GetAchievementByID(ref.MongoAchievementID)
		if err != nil {
			continue
		}

		response := model.AchievementResponse{
			ID:              ref.ID,
			StudentID:       achievement.StudentID,
			AchievementType: achievement.AchievementType,
			Title:           achievement.Title,
			Description:     achievement.Description,
			Details:         achievement.Details,
			Attachments:     achievement.Attachments,
			Tags:            achievement.Tags,
			Points:          achievement.Points,
			Status:          ref.Status,
			CreatedAt:       achievement.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       achievement.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if ref.SubmittedAt != nil {
			submittedAt := ref.SubmittedAt.Format("2006-01-02 15:04:05")
			response.SubmittedAt = &submittedAt
		}

		if ref.VerifiedAt != nil {
			verifiedAt := ref.VerifiedAt.Format("2006-01-02 15:04:05")
			response.VerifiedAt = &verifiedAt
		}

		response.VerifiedBy = ref.VerifiedBy
		response.RejectionNote = ref.RejectionNote

		achievements = append(achievements, response)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: fiber.Map{
			"student_id":   studentID,
			"achievements": achievements,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"total_pages":  totalPages,
		},
	})
}

//
// ==================== SET ADVISOR (PUT /students/:id/advisor) ======================
// SRS Section 5.5: PUT /api/v1/students/:id/advisor
// FR-009: Admin dapat set advisor untuk mahasiswa
//

func (s *StudentService) SetAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// Cari student
	_, err := s.studentRepo.FindByID(studentID)
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

	// Validasi
	if err := s.validate.Struct(req); err != nil {
		return c.Status(422).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	// Validasi advisor exists
	advisor, err := s.lecturerRepo.FindByID(req.AdvisorID)
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

	// Get advisor user info
	advisorUser, _ := s.userRepo.FindByID(advisor.UserID)

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "advisor set successfully",
		Data: fiber.Map{
			"student_id": studentID,
			"advisor": map[string]interface{}{
				"id":          advisor.ID,
				"lecturer_id": advisor.LecturerID,
				"full_name":   advisorUser.FullName,
				"department":  advisor.Department,
			},
		},
	})
}