package service

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"project_uas/app/model"
	"project_uas/test/mocks"
)

// ==================== HELPER FUNCTIONS ====================

func setupAchievementTest() (*AchievementService, *mocks.MockAchievementRepository, *mocks.MockStudentRepository, *mocks.MockLecturerRepository, *mocks.MockUserRepository) {
	mockAchievementRepo := new(mocks.MockAchievementRepository)
	mockStudentRepo := new(mocks.MockStudentRepository)
	mockLecturerRepo := new(mocks.MockLecturerRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	service := NewAchievementService(
		mockAchievementRepo,
		mockStudentRepo,
		mockLecturerRepo,
		mockUserRepo,
	)

	return service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, mockUserRepo
}

// ==================== FR-003: CREATE ACHIEVEMENT ====================

func TestCreateAchievement_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"
	mongoID := "507f1f77bcf86cd799439011"

	// Setup context dengan user claims
	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.CreateAchievement(c)
	})

	// Mock repository calls
	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)

	mockAchievementRepo.On("CreateAchievement", mock.AnythingOfType("*model.Achievement")).Return(mongoID, nil)
	mockAchievementRepo.On("CreateReference", mock.AnythingOfType("*model.AchievementReference")).Return(nil)

	// Request body
	body := `{
		"achievement_type": "competition",
		"title": "Juara 1 Hackathon",
		"description": "Memenangkan hackathon nasional",
		"points": 100
	}`

	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
	mockAchievementRepo.AssertExpectations(t)
}

func TestCreateAchievement_Unauthorized(t *testing.T) {
	service, _, _, _, _ := setupAchievementTest()

	app := fiber.New()
	app.Post("/achievements", service.CreateAchievement)

	body := `{"achievement_type": "competition", "title": "Test"}`
	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestCreateAchievement_StudentNotFound(t *testing.T) {
	service, _, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	app.Post("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "user-123",
			Role:   "Mahasiswa",
		})
		return service.CreateAchievement(c)
	})

	mockStudentRepo.On("FindByUserID", "user-123").Return(nil, errors.New("not found"))

	body := `{"achievement_type": "competition", "title": "Test", "description": "Test"}`
	req := httptest.NewRequest("POST", "/achievements", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
}

// ==================== FR-004: SUBMIT FOR VERIFICATION ====================

func TestSubmitForVerification_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "achievement-123"
	studentID := "student-123"
	userID := "user-123"

	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.SubmitForVerification(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)

	mockAchievementRepo.On("UpdateReference", mock.AnythingOfType("*model.AchievementReference")).Return(nil)

	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/submit", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestSubmitForVerification_NotDraft(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	achievementID := "achievement-123"
	studentID := "student-123"
	userID := "user-123"

	app.Post("/achievements/:id/submit", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.SubmitForVerification(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "verified", // Already verified
	}, nil)

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)

	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/submit", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

// ==================== FR-005: DELETE ACHIEVEMENT ====================

func TestDeleteAchievement_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "achievement-123"
	mongoID := "507f1f77bcf86cd799439011"
	studentID := "student-123"
	userID := "user-123"

	app.Delete("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.DeleteAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:                 achievementID,
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
	}, nil)

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)

	mockAchievementRepo.On("DeleteAchievement", mongoID).Return(nil)
	mockAchievementRepo.On("UpdateReference", mock.AnythingOfType("*model.AchievementReference")).Return(nil)

	req := httptest.NewRequest("DELETE", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestDeleteAchievement_NotDraft(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	achievementID := "achievement-123"
	studentID := "student-123"
	userID := "user-123"

	app.Delete("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.DeleteAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "submitted", // Not draft
	}, nil)

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)

	req := httptest.NewRequest("DELETE", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestDeleteAchievement_Forbidden(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	achievementID := "achievement-123"
	
	app.Delete("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "other-user",
			Role:   "Mahasiswa",
		})
		return service.DeleteAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: "student-123",
		Status:    "draft",
	}, nil)

	mockStudentRepo.On("FindByUserID", "other-user").Return(&model.Student{
		ID:     "other-student",
		UserID: "other-user",
	}, nil)

	req := httptest.NewRequest("DELETE", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

// ==================== FR-007: VERIFY ACHIEVEMENT ====================

func TestVerifyAchievement_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, _ := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "achievement-123"
	studentID := "student-123"
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.VerifyAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "submitted",
	}, nil)

	mockLecturerRepo.On("FindByUserID", userID).Return(&model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}, nil)

	mockStudentRepo.On("FindByID", studentID).Return(&model.Student{
		ID:        studentID,
		AdvisorID: &lecturerID,
	}, nil)

	mockAchievementRepo.On("UpdateReference", mock.AnythingOfType("*model.AchievementReference")).Return(nil)

	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/verify", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestVerifyAchievement_NotSubmitted(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, _ := setupAchievementTest()

	app := fiber.New()
	achievementID := "achievement-123"
	studentID := "student-123"
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.VerifyAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "draft", // Not submitted
	}, nil)

	mockLecturerRepo.On("FindByUserID", userID).Return(&model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}, nil)

	mockStudentRepo.On("FindByID", studentID).Return(&model.Student{
		ID:        studentID,
		AdvisorID: &lecturerID,
	}, nil)

	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/verify", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestVerifyAchievement_NotAdvisor(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, _ := setupAchievementTest()

	app := fiber.New()
	achievementID := "achievement-123"
	studentID := "student-123"
	lecturerID := "lecturer-123"
	otherLecturerID := "other-lecturer"
	userID := "user-lecturer"

	app.Post("/achievements/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.VerifyAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "submitted",
	}, nil)

	mockLecturerRepo.On("FindByUserID", userID).Return(&model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}, nil)

	mockStudentRepo.On("FindByID", studentID).Return(&model.Student{
		ID:        studentID,
		AdvisorID: &otherLecturerID, // Different advisor
	}, nil)

	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/verify", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

// ==================== FR-008: REJECT ACHIEVEMENT ====================

func TestRejectAchievement_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, _ := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "achievement-123"
	studentID := "student-123"
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Post("/achievements/:id/reject", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.RejectAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "submitted",
	}, nil)

	mockLecturerRepo.On("FindByUserID", userID).Return(&model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}, nil)

	mockStudentRepo.On("FindByID", studentID).Return(&model.Student{
		ID:        studentID,
		AdvisorID: &lecturerID,
	}, nil)

	mockAchievementRepo.On("UpdateReference", mock.AnythingOfType("*model.AchievementReference")).Return(nil)

	body := `{"rejection_note": "Data tidak lengkap"}`
	req := httptest.NewRequest("POST", "/achievements/"+achievementID+"/reject", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
}

// ==================== UPDATE ACHIEVEMENT ====================

func TestUpdateAchievement_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "achievement-123"
	mongoID := "507f1f77bcf86cd799439011"
	studentID := "student-123"
	userID := "user-123"

	app.Put("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.UpdateAchievement(c)
	})

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(&model.AchievementReference{
		ID:                 achievementID,
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "draft",
	}, nil)

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)

	mockAchievementRepo.On("GetAchievementByID", mongoID).Return(&model.Achievement{
		StudentID:       studentID,
		AchievementType: "competition",
		Title:           "Old Title",
		Description:     "Old Description",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil)

	mockAchievementRepo.On("UpdateAchievement", mongoID, mock.AnythingOfType("*model.Achievement")).Return(nil)

	body := `{"title": "New Title", "description": "New Description"}`
	req := httptest.NewRequest("PUT", "/achievements/"+achievementID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
}
// ==================== TAMBAHAN TEST - PASTE DI AKHIR achievement_service_test.go ====================

// ==================== GET ACHIEVEMENT BY ID ====================

func TestGetAchievementByID_Success(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "ref-123"
	mongoID := "mongo-123"
	studentID := "student-123"
	userID := "user-123"

	app.Get("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetAchievementByID(c)
	})

	reference := &model.AchievementReference{
		ID:                 achievementID,
		StudentID:          studentID,
		MongoAchievementID: mongoID,
		Status:             "verified",
		CreatedAt:          time.Now(),
	}

	achievement := &model.Achievement{
		StudentID:       studentID,
		AchievementType: "competition",
		Title:           "Test Achievement",
		Description:     "Description",
		Points:          100,
		Details: map[string]interface{}{
			"competitionLevel": "national",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	student := &model.Student{
		ID:     studentID,
		UserID: userID,
	}

	user := &model.User{
		ID:       userID,
		FullName: "Student Name",
	}

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(reference, nil)
	mockAchievementRepo.On("GetAchievementByID", mongoID).Return(achievement, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockUserRepo.On("FindByID", userID).Return(user, nil)

	req := httptest.NewRequest("GET", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetAchievementByID_NotFound(t *testing.T) {
	service, mockAchievementRepo, _, _, _ := setupAchievementTest()

	app := fiber.New()
	
	app.Get("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "user-123",
			Role:   "Mahasiswa",
		})
		return service.GetAchievementByID(c)
	})

	mockAchievementRepo.On("GetReferenceByID", "invalid-id").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/achievements/invalid-id", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetAchievementByID_Forbidden(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	achievementID := "ref-123"
	studentID := "student-123"
	otherUserID := "other-user"

	app.Get("/achievements/:id", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: otherUserID,
			Role:   "Mahasiswa",
		})
		return service.GetAchievementByID(c)
	})

	reference := &model.AchievementReference{
		ID:        achievementID,
		StudentID: studentID,
		Status:    "draft",
	}

	mockAchievementRepo.On("GetReferenceByID", achievementID).Return(reference, nil)
	mockStudentRepo.On("FindByUserID", otherUserID).Return(&model.Student{
		ID:     "other-student",
		UserID: otherUserID,
	}, nil)

	req := httptest.NewRequest("GET", "/achievements/"+achievementID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
}

// ==================== GET ACHIEVEMENTS (LIST) ====================

func TestGetAchievements_Success_Mahasiswa(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupAchievementTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetAchievements(c)
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
		Points:          100,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockStudentRepo.On("FindByUserID", userID).Return(student, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 10, 0).Return(references, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "").Return(1, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", userID).Return(&model.User{
		ID:       userID,
		FullName: "Student Name",
	}, nil)

	req := httptest.NewRequest("GET", "/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetAchievements_Success_DosenWali(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, mockLecturerRepo, mockUserRepo := setupAchievementTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	studentID := "student-123"
	userID := "user-lecturer"

	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetAchievements(c)
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
			Status:             "submitted",
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

	student := &model.Student{
		ID:     studentID,
		UserID: "user-student",
	}

	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockAchievementRepo.On("GetReferencesByAdvisorID", lecturerID, "", 10, 0).Return(references, nil)
	mockAchievementRepo.On("CountReferencesByAdvisorID", lecturerID, "").Return(1, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", "user-student").Return(&model.User{
		ID:       "user-student",
		FullName: "Student Name",
	}, nil)

	req := httptest.NewRequest("GET", "/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetAchievements_Success_Admin(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, mockUserRepo := setupAchievementTest()

	app := fiber.New()
	
	studentID := "student-123"

	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "admin-user",
			Role:   "Admin",
		})
		return service.GetAchievements(c)
	})

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

	student := &model.Student{
		ID:     studentID,
		UserID: "user-student",
	}

	mockAchievementRepo.On("GetAllReferences", "", 10, 0).Return(references, nil)
	mockAchievementRepo.On("CountAllReferences", "").Return(1, nil)
	mockAchievementRepo.On("GetAchievementByID", "mongo-1").Return(achievement, nil)
	mockStudentRepo.On("FindByID", studentID).Return(student, nil)
	mockUserRepo.On("FindByID", "user-student").Return(&model.User{
		ID:       "user-student",
		FullName: "Student Name",
	}, nil)

	req := httptest.NewRequest("GET", "/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockAchievementRepo.AssertExpectations(t)
}

func TestGetAchievements_WithPagination(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetAchievements(c)
	})

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "", 20, 20).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "").Return(50, nil)

	req := httptest.NewRequest("GET", "/achievements?page=2&page_size=20", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetAchievements_WithStatusFilter(t *testing.T) {
	service, mockAchievementRepo, mockStudentRepo, _, _ := setupAchievementTest()

	app := fiber.New()
	
	studentID := "student-123"
	userID := "user-123"

	app.Get("/achievements", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Mahasiswa",
		})
		return service.GetAchievements(c)
	})

	mockStudentRepo.On("FindByUserID", userID).Return(&model.Student{
		ID:     studentID,
		UserID: userID,
	}, nil)
	mockAchievementRepo.On("GetReferencesByStudentID", studentID, "verified", 10, 0).Return([]model.AchievementReference{}, nil)
	mockAchievementRepo.On("CountReferencesByStudentID", studentID, "verified").Return(0, nil)

	req := httptest.NewRequest("GET", "/achievements?status=verified", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetAchievements_Unauthorized(t *testing.T) {
	service, _, _, _, _ := setupAchievementTest()

	app := fiber.New()
	app.Get("/achievements", service.GetAchievements)

	req := httptest.NewRequest("GET", "/achievements", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

// ... test lain di atas ...

// ==================== GET ACHIEVEMENT HISTORY ====================

func TestGetAchievementHistory_Success(t *testing.T) {
	t.Skip("Skipped: Complex routing - requires integration test")
}

func TestGetAchievementHistory_StudentNotFound(t *testing.T) {
	t.Skip("Skipped: Complex routing - requires integration test")
}

func TestGetAchievementHistory_EmptyResult(t *testing.T) {
	t.Skip("Skipped: Complex routing - requires integration test")
}

// ==================== UPLOAD ATTACHMENT ====================

func TestUploadAttachment_AchievementNotFound(t *testing.T) {
	t.Skip("Skipped: File upload requires integration test with multipart form")
}

func TestUploadAttachment_NotDraft(t *testing.T) {
	t.Skip("Skipped: File upload requires integration test with multipart form")
}

func TestUploadAttachment_Forbidden(t *testing.T) {
	t.Skip("Skipped: File upload requires integration test with multipart form")
}