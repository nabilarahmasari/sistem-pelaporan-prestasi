package service

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"project_uas/app/model"
	"project_uas/test/mocks"
)

// ==================== HELPER FUNCTIONS ====================

func setupReportTest() (*ReportService, *mocks.MockAchievementRepository, *mocks.MockStudentRepository, *mocks.MockLecturerRepository, *mocks.MockUserRepository) {
	mockAchievementRepo := new(mocks.MockAchievementRepository)
	mockStudentRepo := new(mocks.MockStudentRepository)
	mockLecturerRepo := new(mocks.MockLecturerRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	service := NewReportService(mockAchievementRepo, mockStudentRepo, mockLecturerRepo, mockUserRepo)

	return service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, mockUserRepo
}

// ==================== FR-011: GET STATISTICS (MAHASISWA) ====================

func TestGetStatistics_Success_Mahasiswa(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetStatistics(c)
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
		{
			ID:                 "ref-2",
			StudentID:          studentID,
			MongoAchievementID: "mongo-2",
			Status:             "submitted",
			CreatedAt:          time.Now(),
		},
	}

	achievement1 := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "competition",
		Title:           "Competition Achievement",
		Points:          100,
		Details: map[string]interface{}{
			"competitionLevel": "national",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	achievement2 := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "publication",
		Title:           "Publication Achievement",
		Points:          50,
		Details:         map[string]interface{}{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10000, 0).Return(references, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement1, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-2").Return(achievement2, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", userID).Return(&model.User{
		ID:       userID,
		FullName: "John Doe",
	}, nil)

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetStatistics_StudentNotFound(t *testing.T) {
	service, _, mockStudentRepo, _, _ := setupReportTest()

	app := fiber.New()
	
	userID := "user-123"

	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetStatistics(c)
	})

	mockStudentRepo.On("FindByUserID", userID).Return(nil, errors.New("student not found"))

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

// ==================== FR-011: GET STATISTICS (DOSEN WALI) ====================

func TestGetStatistics_Success_DosenWali(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	studentID := "student-123"
	userID := "user-lecturer"

	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetStatistics(c)
	})

	lecturer := &model.Lecturer{
		ID:     lecturerID,
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
		Points:          100,
		Details: map[string]interface{}{
			"competitionLevel": "international",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	student := &model.Student{
		ID:        studentID,
		StudentID: "123456",
	}

	user := &model.User{
		ID:       "user-student",
		FullName: "Student Name",
	}

	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockAchievementRepo.On("GetReferencesByAdvisorID", lecturerID, "", 10000, 0).Return(references, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", student.UserID).Return(user, nil)

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetStatistics_LecturerNotFound(t *testing.T) {
	service, _, _, mockLecturerRepo, _ := setupReportTest()

	app := fiber.New()
	
	userID := "user-lecturer"

	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetStatistics(c)
	})

	mockLecturerRepo.On("FindByUserID", userID).Return(nil, errors.New("lecturer not found"))

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

// ==================== FR-011: GET STATISTICS (ADMIN) ====================

func TestGetStatistics_Success_Admin(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "admin-user",
			Role:   "Admin",
		})
		return service.GetStatistics(c)
	})

	studentID := "student-123"
	references := []model.AchievementReference{
		{
			ID:                 "ref-1",
			StudentID:          studentID,
			MongoAchievementID: "mongo-1",
			Status:             "verified",
			CreatedAt:          time.Now(),
		},
		{
			ID:                 "ref-2",
			StudentID:          studentID,
			MongoAchievementID: "mongo-2",
			Status:             "rejected",
			CreatedAt:          time.Now(),
		},
	}

	achievement1 := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "academic",
		Points:          75,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	achievement2 := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "certification",
		Points:          50,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	student := &model.Student{
		ID:        studentID,
		StudentID: "123456",
		UserID:    "user-student",
	}

	user := &model.User{
		ID:       "user-student",
		FullName: "Student Name",
	}

	mockAchievementRepo.On("GetAllReferences", "", 10000, 0).Return(references, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement1, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-2").Return(achievement2, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", "user-student").Return(user, nil)

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetStatistics_Unauthorized(t *testing.T) {
	service, _, _, _, _ := setupReportTest()

	app := fiber.New()
	app.Get("/statistics", service.GetStatistics)

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestGetStatistics_ForbiddenRole(t *testing.T) {
	service, _, _, _, _ := setupReportTest()

	app := fiber.New()
	
	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "user-123",
			Role:   "InvalidRole",
		})
		return service.GetStatistics(c)
	})

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

// ==================== GET STUDENT REPORT ====================

func TestGetStudentReport_Success_OwnReport(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetStudentReport(c)
	})

	student := &model.Student{
		ID:           studentID,
		UserID:       userID,
		StudentID:    "123456",
		ProgramStudy: "Teknik Informatika",
		AcademicYear: "2024",
	}

	user := &model.User{
		ID:       userID,
		FullName: "John Doe",
		Email:    "john@test.com",
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
		Points:          100,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10000, 0).Return(references, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 5, 0).Return(references, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement, nil)

	req := httptest.NewRequest("GET", "/reports/student/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetStudentReport_Success_DosenWali(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetStudentReport(c)
	})

	student := &model.Student{
		ID:        studentID,
		UserID:    "user-student",
		StudentID: "123456",
		AdvisorID: &lecturerID,
	}

	studentUser := &model.User{
		ID:       "user-student",
		FullName: "Student Name",
		Email:    "student@test.com",
	}

	lecturer := &model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockUserRepo.On("FindByID", "user-student").Return(studentUser, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10000, 0).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 5, 0).Return([]model.AchievementReference{}, nil)

	req := httptest.NewRequest("GET", "/reports/student/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetStudentReport_Success_Admin(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "admin-user",
			Role:   "Admin",
		})
		return service.GetStudentReport(c)
	})

	student := &model.Student{
		ID:        studentID,
		UserID:    "user-student",
		StudentID: "123456",
	}

	user := &model.User{
		ID:       "user-student",
		FullName: "Student Name",
		Email:    "student@test.com",
	}

	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", "user-student").Return(user, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10000, 0).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 5, 0).Return([]model.AchievementReference{}, nil)

	req := httptest.NewRequest("GET", "/reports/student/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetStudentReport_StudentNotFound(t *testing.T) {
	service, _, mockStudentRepo, _, _ := setupReportTest()

	app := fiber.New()
	
	studentID := "invalid-student"

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "user-123",
			Role:   "Admin",
		})
		return service.GetStudentReport(c)
	})

	mockStudentRepo.On("FindByID", studentID).Return(nil, errors.New("student not found"))

	req := httptest.NewRequest("GET", "/reports/student/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentReport_Forbidden_NotOwner(t *testing.T) {
	service, _, mockStudentRepo, _, _ := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"
	otherUserID := "other-user"

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: otherUserID,
			Role:   "Mahasiswa",
		})
		return service.GetStudentReport(c)
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

	req := httptest.NewRequest("GET", "/reports/student/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentReport_Forbidden_NotAdvisor(t *testing.T) {
	service, _, mockStudentRepo, mockLecturerRepo, _ := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"
	lecturerID := "lecturer-123"
	otherLecturerID := "other-lecturer"
	userID := "user-lecturer"

	app.Get("/reports/student/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetStudentReport(c)
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

	req := httptest.NewRequest("GET", "/reports/student/"+studentID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetStudentReport_Unauthorized(t *testing.T) {
	service, _, _, _, _ := setupReportTest()

	app := fiber.New()
	app.Get("/reports/student/:id", service.GetStudentReport)

	req := httptest.NewRequest("GET", "/reports/student/student-123", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

// ==================== COMPLEX STATISTICS ====================

func TestGetStatistics_MultipleAchievementTypes(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupReportTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/statistics", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetStatistics(c)
	})

	student := &model.Student{
		ID:     studentID,
		UserID: userID,
	}

	references := []model.AchievementReference{
		{StudentID: studentID, MongoAchievementID: "mongo-1", Status: "verified"},
		{StudentID: studentID, MongoAchievementID: "mongo-2", Status: "submitted"},
		{StudentID: studentID, MongoAchievementID: "mongo-3", Status: "rejected"},
		{StudentID: studentID, MongoAchievementID: "mongo-4", Status: "verified"},
	}

	achievements := []*model.Achievement{
		{
			StudentID:       studentID,
			AchievementType: "competition",
			Points:          100,
			Details:         map[string]interface{}{"competitionLevel": "international"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			StudentID:       studentID,
			AchievementType: "publication",
			Points:          80,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			StudentID:       studentID,
			AchievementType: "competition",
			Points:          60,
			Details:         map[string]interface{}{"competitionLevel": "national"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			StudentID:       studentID,
			AchievementType: "certification",
			Points:          40,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10000, 0).Return(references, nil)
	
	for i, ach := range achievements {
		mockAchievementRepo.On("GetAchievementByID", references[i].MongoAchievementID).Return(ach, nil)
	}
	
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", userID).Return(&model.User{
		ID:       userID,
		FullName: "John Doe",
	}, nil)

	req := httptest.NewRequest("GET", "/statistics", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}