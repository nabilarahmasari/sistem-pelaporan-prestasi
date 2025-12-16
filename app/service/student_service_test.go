package service

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"


	"project_uas/app/model"
	"project_uas/test/mocks"
)

// ==================== HELPER FUNCTIONS ====================

func setupStudentTest() (*StudentService, *mocks.MockStudentRepository, *mocks.MockLecturerRepository, *mocks.MockUserRepository, *mocks.MockAchievementRepository) {
	mockStudentRepo := new(mocks.MockStudentRepository)
	mockLecturerRepo := new(mocks.MockLecturerRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockAchievementRepo := new(mocks.MockAchievementRepository)

	service := NewStudentService(mockStudentRepo, mockLecturerRepo, mockUserRepo, mockAchievementRepo)

	return service, mockStudentRepo, mockLecturerRepo, mockUserRepo, mockAchievementRepo
}

// ==================== GET ALL STUDENTS ====================

func TestGetAllStudents_Success(t *testing.T) {
	service, mockStudentRepo, _, mockUserRepo, _ := setupStudentTest()

	app := fiber.New()
	app.Get("/students", service.GetAllStudents)

	students := []model.Student{
		{
			ID:           "student-1",
			UserID:       "user-1",
			StudentID:    "123456",
			ProgramStudy: "Teknik Informatika",
			AcademicYear: "2024",
			CreatedAt:    time.Now(),
		},
		{
			ID:           "student-2",
			UserID:       "user-2",
			StudentID:    "123457",
			ProgramStudy: "Sistem Informasi",
			AcademicYear: "2024",
			CreatedAt:    time.Now(),
		},
	}

	user1 := &model.User{
		ID:       "user-1",
		Username: "mhs1",
		FullName: "Student One",
		Email:    "student1@test.com",
	}

	user2 := &model.User{
		ID:       "user-2",
		Username: "mhs2",
		FullName: "Student Two",
		Email:    "student2@test.com",
	}

	mockStudentRepo.On("GetAll", 10, 0).Return(students, nil)
	mockStudentRepo.On("CountAll").Return(2, nil)
	mockUserRepo.On("FindByID", "user-1").Return(user1, nil)
	mockUserRepo.On("FindByID", "user-2").Return(user2, nil)

	req := httptest.NewRequest("GET", "/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetAllStudents_WithPagination(t *testing.T) {
	service, mockStudentRepo, _, _, _ := setupStudentTest()

	app := fiber.New()
	app.Get("/students", service.GetAllStudents)

	mockStudentRepo.On("GetAll", 20, 20).Return([]model.Student{}, nil)
	mockStudentRepo.On("CountAll").Return(50, nil)

	req := httptest.NewRequest("GET", "/students?page=2&page_size=20", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestGetAllStudents_EmptyResult(t *testing.T) {
	service, mockStudentRepo, _, _, _ := setupStudentTest()

	app := fiber.New()
	app.Get("/students", service.GetAllStudents)

	mockStudentRepo.On("GetAll", 10, 0).Return([]model.Student{}, nil)
	mockStudentRepo.On("CountAll").Return(0, nil)

	req := httptest.NewRequest("GET", "/students", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

// ==================== GET STUDENT BY ID ====================

func TestGetStudentByID_Success(t *testing.T) {
	service, mockStudentRepo, mockLecturerRepo, mockUserRepo, _ := setupStudentTest()

	app := fiber.New()
	app.Get("/students/:id", service.GetStudentByID)

	studentID := "student-123"
	advisorID := "advisor-123"

	student := &model.Student{
		ID:           studentID,
		UserID:       "user-123",
		StudentID:    "123456",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2024",
		AdvisorID:    &advisorID,
		CreatedAt:    time.Now(),
	}

	user := &model.User{
		ID:       "user-123",
		Username: "mahasiswa123",
		FullName: "John Doe",
		Email:    "john@test.com",
	}

	lecturer := &model.Lecturer{
		ID:         advisorID,
		UserID:     "user-lecturer",
		LecturerID: "L123",
		Department: "Informatika",
	}

	lecturerUser := &model.User{
		ID:       "user-lecturer",
		FullName: "Dr. Dosen",
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", "user-123").Return(user, nil)
	mockLecturerRepo.On("FindByID", advisorID).Return(lecturer, nil)
	mockUserRepo.On("FindByID", "user-lecturer").Return(lecturerUser, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetStudentByID_NotFound(t *testing.T) {
	service, mockStudentRepo, _, _, _ := setupStudentTest()

	app := fiber.New()
	app.Get("/students/:id", service.GetStudentByID)

	mockStudentRepo.On("FindByID", "invalid-id").Return(nil, errors.New("student not found"))

	req := httptest.NewRequest("GET", "/students/invalid-id", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

// ==================== GET STUDENT ACHIEVEMENTS ====================

func TestGetStudentAchievements_Success_Mahasiswa(t *testing.T) {
	service, mockStudentRepo, _, _, mockAchievementRepo := setupStudentTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetStudentAchievements(c)
	})

	student := &model.Student{
		ID:     studentID,
		UserID: userID,
	}

	references := []model.AchievementReference{
		{
			ID:                 "ref-1",
			StudentID:          studentID,
			MongoAchievementID: "mongo-1",
			Status:             "verified",
			CreatedAt:          time.Now(),
		},
	}

	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "competition",
		Title:           "Test Achievement",
		Description:     "Description",
		Points:          100,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10, 0).Return(references, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "").Return(1, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID+"/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetStudentAchievements_Success_DosenWali(t *testing.T) {
	service, mockStudentRepo, mockLecturerRepo, _, mockAchievementRepo := setupStudentTest()

	app := fiber.New()
	
	studentID := "student-123"
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetStudentAchievements(c)
	})

	student := &model.Student{
		ID:        studentID,
		UserID:    "user-student",
		AdvisorID: &lecturerID,
	}

	lecturer := &model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10, 0).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "").Return(0, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID+"/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetStudentAchievements_Forbidden_NotOwner(t *testing.T) {
	service, mockStudentRepo, _, _, _ := setupStudentTest()

	app := fiber.New()
	
	studentID := "student-123"
	otherUserID := "other-user"

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: otherUserID,
			Role:   "Mahasiswa",
		})
		return service.GetStudentAchievements(c)
	})

	student := &model.Student{
		ID:     studentID,
		UserID: "user-123", // Different user
	}

	otherStudent := &model.Student{
		ID:     "other-student",
		UserID: otherUserID,
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockStudentRepo.On("FindByUserID", otherUserID).Return(otherStudent, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID+"/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentAchievements_Forbidden_NotAdvisor(t *testing.T) {
	service, mockStudentRepo, mockLecturerRepo, _, _ := setupStudentTest()

	app := fiber.New()
	
	studentID := "student-123"
	lecturerID := "lecturer-123"
	otherLecturerID := "other-lecturer"
	userID := "user-lecturer"

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetStudentAchievements(c)
	})

	student := &model.Student{
		ID:        studentID,
		UserID:    "user-student",
		AdvisorID: &otherLecturerID, // Different advisor
	}

	lecturer := &model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID+"/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetStudentAchievements_Success_Admin(t *testing.T) {
	service, mockStudentRepo, _, _, mockAchievementRepo := setupStudentTest()

	app := fiber.New()
	
	studentID := "student-123"

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "admin-user",
			Role:   "Admin",
		})
		return service.GetStudentAchievements(c)
	})

	student := &model.Student{
		ID:     studentID,
		UserID: "user-123",
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10, 0).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "").Return(0, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID+"/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetStudentAchievements_WithStatusFilter(t *testing.T) {
	service, mockStudentRepo, _, _, mockAchievementRepo := setupStudentTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/students/:id/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetStudentAchievements(c)
	})

	student := &model.Student{
		ID:     studentID,
		UserID: userID,
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "verified", 10, 0).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "verified").Return(0, nil)

	req := httptest.NewRequest("GET", "/students/"+studentID+"/achievements?status=verified", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

// ==================== SET ADVISOR ====================

func TestSetAdvisor_Success(t *testing.T) {
	service, mockStudentRepo, mockLecturerRepo, mockUserRepo, _ := setupStudentTest()

	app := fiber.New()
	app.Put("/students/:id/advisor", service.SetAdvisor)

	studentID := "student-123"
	advisorID := "advisor-123"

	student := &model.Student{
		ID:        studentID,
		UserID:    "user-123",
		StudentID: "123456",
	}

	advisor := &model.Lecturer{
		ID:         advisorID,
		UserID:     "user-advisor",
		LecturerID: "L123",
		Department: "Informatika",
	}

	advisorUser := &model.User{
		ID:       "user-advisor",
		FullName: "Dr. Dosen",
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockLecturerRepo.On("FindByID", advisorID).Return(advisor, nil)
	mockStudentRepo.On("SetAdvisor", studentID, advisorID).Return(nil)
	mockUserRepo.On("FindByID", "user-advisor").Return(advisorUser, nil)

	body := `{"advisor_id": "` + advisorID + `"}`

	req := httptest.NewRequest("PUT", "/students/"+studentID+"/advisor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestSetAdvisor_StudentNotFound(t *testing.T) {
	service, mockStudentRepo, _, _, _ := setupStudentTest()

	app := fiber.New()
	app.Put("/students/:id/advisor", service.SetAdvisor)

	studentID := "invalid-student"

	mockStudentRepo.On("FindByID", studentID).Return(nil, errors.New("student not found"))

	body := `{"advisor_id": "advisor-123"}`

	req := httptest.NewRequest("PUT", "/students/"+studentID+"/advisor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestSetAdvisor_AdvisorNotFound(t *testing.T) {
	service, mockStudentRepo, mockLecturerRepo, _, _ := setupStudentTest()

	app := fiber.New()
	app.Put("/students/:id/advisor", service.SetAdvisor)

	studentID := "student-123"
	advisorID := "invalid-advisor"

	student := &model.Student{
		ID:     studentID,
		UserID: "user-123",
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockLecturerRepo.On("FindByID", advisorID).Return(nil, errors.New("lecturer not found"))

	body := `{"advisor_id": "` + advisorID + `"}`

	req := httptest.NewRequest("PUT", "/students/"+studentID+"/advisor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
}

func TestSetAdvisor_InvalidRequestBody(t *testing.T) {
	service, mockStudentRepo, _, _, _ := setupStudentTest()

	app := fiber.New()
	app.Put("/students/:id/advisor", service.SetAdvisor)

	studentID := "student-123"

	student := &model.Student{
		ID:     studentID,
		UserID: "user-123",
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)

	body := `{"invalid": "data"}`

	req := httptest.NewRequest("PUT", "/students/"+studentID+"/advisor", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 422, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}