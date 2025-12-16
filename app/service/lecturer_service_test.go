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

func setupLecturerTest() (*LecturerService, *mocks.MockLecturerRepository, *mocks.MockStudentRepository, *mocks.MockUserRepository) {
	mockLecturerRepo := new(mocks.MockLecturerRepository)
	mockStudentRepo := new(mocks.MockStudentRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	service := NewLecturerService(mockLecturerRepo, mockStudentRepo, mockUserRepo)

	return service, mockLecturerRepo, mockStudentRepo, mockUserRepo
}
// ==================== GET ALL LECTURERS ====================

func TestGetAllLecturers_Success(t *testing.T) {
	service, mockLecturerRepo, _, mockUserRepo := setupLecturerTest()

	app := fiber.New()
	app.Get("/lecturers", service.GetAllLecturers)

	lecturers := []model.Lecturer{
		{
			ID:         "lecturer-1",
			UserID:     "user-1",
			LecturerID: "L001",
			Department: "Informatika",
			CreatedAt:  time.Now(),
		},
		{
			ID:         "lecturer-2",
			UserID:     "user-2",
			LecturerID: "L002",
			Department: "Sistem Informasi",
			CreatedAt:  time.Now(),
		},
	}

	user1 := &model.User{
		ID:       "user-1",
		Username: "dosen1",
		FullName: "Dr. Dosen Satu",
		Email:    "dosen1@test.com",
	}

	user2 := &model.User{
		ID:       "user-2",
		Username: "dosen2",
		FullName: "Dr. Dosen Dua",
		Email:    "dosen2@test.com",
	}

	mockLecturerRepo.On("GetAll", 10, 0).Return(lecturers, nil)
	mockLecturerRepo.On("CountAll").Return(2, nil)
	mockUserRepo.On("FindByID", "user-1").Return(user1, nil)
	mockUserRepo.On("FindByID", "user-2").Return(user2, nil)

	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetAllLecturers_WithPagination(t *testing.T) {
	service, mockLecturerRepo, _, _ := setupLecturerTest()

	app := fiber.New()
	app.Get("/lecturers", service.GetAllLecturers)

	mockLecturerRepo.On("GetAll", 20, 20).Return([]model.Lecturer{}, nil)
	mockLecturerRepo.On("CountAll").Return(50, nil)

	req := httptest.NewRequest("GET", "/lecturers?page=2&page_size=20", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetAllLecturers_EmptyResult(t *testing.T) {
	service, mockLecturerRepo, _, _ := setupLecturerTest()

	app := fiber.New()
	app.Get("/lecturers", service.GetAllLecturers)

	mockLecturerRepo.On("GetAll", 10, 0).Return([]model.Lecturer{}, nil)
	mockLecturerRepo.On("CountAll").Return(0, nil)

	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetAllLecturers_RepositoryError(t *testing.T) {
	service, mockLecturerRepo, _, _ := setupLecturerTest()

	app := fiber.New()
	app.Get("/lecturers", service.GetAllLecturers)

	mockLecturerRepo.On("GetAll", 10, 0).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/lecturers", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

// ==================== FR-006: GET LECTURER ADVISEES ====================

func TestGetLecturerAdvisees_Success_OwnAdvisees(t *testing.T) {
	service, mockLecturerRepo, mockStudentRepo, mockUserRepo := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetLecturerAdvisees(c)
	})

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "L123",
		Department: "Informatika",
	}

	lecturerUser := &model.User{
		ID:       userID,
		FullName: "Dr. Dosen",
	}

	allStudents := []model.Student{
		{
			ID:           "student-1",
			UserID:       "user-student-1",
			StudentID:    "123456",
			ProgramStudy: "Teknik Informatika",
			AdvisorID:    &lecturerID, // Advisee
			CreatedAt:    time.Now(),
		},
		{
			ID:           "student-2",
			UserID:       "user-student-2",
			StudentID:    "123457",
			ProgramStudy: "Sistem Informasi",
			AdvisorID:    nil, // No advisor
			CreatedAt:    time.Now(),
		},
	}

	student1User := &model.User{
		ID:       "user-student-1",
		FullName: "Student One",
		Email:    "student1@test.com",
	}

	mockLecturerRepo.On("FindByID", lecturerID).Return(lecturer, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockStudentRepo.On("GetAll", 1000, 0).Return(allStudents, nil)
	mockUserRepo.On("FindByID", userID).Return(lecturerUser, nil)
	mockUserRepo.On("FindByID", "user-student-1").Return(student1User, nil)

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetLecturerAdvisees_Success_Admin(t *testing.T) {
	service, mockLecturerRepo, mockStudentRepo, mockUserRepo := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "admin-user",
			Role:   "Admin",
		})
		return service.GetLecturerAdvisees(c)
	})

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     "user-lecturer",
		LecturerID: "L123",
		Department: "Informatika",
	}

	lecturerUser := &model.User{
		ID:       "user-lecturer",
		FullName: "Dr. Dosen",
	}

	allStudents := []model.Student{
		{
			ID:        "student-1",
			UserID:    "user-student-1",
			StudentID: "123456",
			AdvisorID: &lecturerID,
			CreatedAt: time.Now(),
		},
	}

	student1User := &model.User{
		ID:       "user-student-1",
		FullName: "Student One",
		Email:    "student1@test.com",
	}

	mockLecturerRepo.On("FindByID", lecturerID).Return(lecturer, nil)
	mockStudentRepo.On("GetAll", 1000, 0).Return(allStudents, nil)
	mockUserRepo.On("FindByID", "user-lecturer").Return(lecturerUser, nil)
	mockUserRepo.On("FindByID", "user-student-1").Return(student1User, nil)

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetLecturerAdvisees_LecturerNotFound(t *testing.T) {
	service, mockLecturerRepo, _, _ := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "invalid-lecturer"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: "user-123",
			Role:   "Admin",
		})
		return service.GetLecturerAdvisees(c)
	})

	mockLecturerRepo.On("FindByID", lecturerID).Return(nil, errors.New("lecturer not found"))

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetLecturerAdvisees_Forbidden_NotOwnAdvisees(t *testing.T) {
	service, mockLecturerRepo, _, _ := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	otherLecturerID := "other-lecturer"
	userID := "user-lecturer"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetLecturerAdvisees(c)
	})

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     "other-user",
		LecturerID: "L123",
	}

	currentLecturer := &model.Lecturer{
		ID:     otherLecturerID, // Different lecturer
		UserID: userID,
	}

	mockLecturerRepo.On("FindByID", lecturerID).Return(lecturer, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(currentLecturer, nil)

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 403, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
}

func TestGetLecturerAdvisees_NoAdvisees(t *testing.T) {
	service, mockLecturerRepo, mockStudentRepo, mockUserRepo := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetLecturerAdvisees(c)
	})

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "L123",
		Department: "Informatika",
	}

	lecturerUser := &model.User{
		ID:       userID,
		FullName: "Dr. Dosen",
	}

	// All students have different advisors
	otherLecturerID := "other-lecturer"
	allStudents := []model.Student{
		{
			ID:        "student-1",
			AdvisorID: &otherLecturerID,
		},
		{
			ID:        "student-2",
			AdvisorID: nil,
		},
	}

	mockLecturerRepo.On("FindByID", lecturerID).Return(lecturer, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockStudentRepo.On("GetAll", 1000, 0).Return(allStudents, nil)
	mockUserRepo.On("FindByID", userID).Return(lecturerUser, nil)

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetLecturerAdvisees_Unauthorized(t *testing.T) {
	service, _, _, _ := setupLecturerTest()

	app := fiber.New()
	app.Get("/lecturers/:id/advisees", service.GetLecturerAdvisees)

	req := httptest.NewRequest("GET", "/lecturers/lecturer-123/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestGetLecturerAdvisees_MultipleAdvisees(t *testing.T) {
	service, mockLecturerRepo, mockStudentRepo, mockUserRepo := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Dosen Wali",
		})
		return service.GetLecturerAdvisees(c)
	})

	lecturer := &model.Lecturer{
		ID:         lecturerID,
		UserID:     userID,
		LecturerID: "L123",
		Department: "Informatika",
	}

	lecturerUser := &model.User{
		ID:       userID,
		FullName: "Dr. Dosen",
	}

	allStudents := []model.Student{
		{
			ID:           "student-1",
			UserID:       "user-1",
			StudentID:    "111111",
			ProgramStudy: "TI",
			AdvisorID:    &lecturerID,
			CreatedAt:    time.Now(),
		},
		{
			ID:           "student-2",
			UserID:       "user-2",
			StudentID:    "222222",
			ProgramStudy: "SI",
			AdvisorID:    &lecturerID,
			CreatedAt:    time.Now(),
		},
		{
			ID:           "student-3",
			UserID:       "user-3",
			StudentID:    "333333",
			ProgramStudy: "TI",
			AdvisorID:    &lecturerID,
			CreatedAt:    time.Now(),
		},
	}

	user1 := &model.User{ID: "user-1", FullName: "Student 1", Email: "s1@test.com"}
	user2 := &model.User{ID: "user-2", FullName: "Student 2", Email: "s2@test.com"}
	user3 := &model.User{ID: "user-3", FullName: "Student 3", Email: "s3@test.com"}

	mockLecturerRepo.On("FindByID", lecturerID).Return(lecturer, nil)
	mockLecturerRepo.On("FindByUserID", userID).Return(lecturer, nil)
	mockStudentRepo.On("GetAll", 1000, 0).Return(allStudents, nil)
	mockUserRepo.On("FindByID", userID).Return(lecturerUser, nil)
	mockUserRepo.On("FindByID", "user-1").Return(user1, nil)
	mockUserRepo.On("FindByID", "user-2").Return(user2, nil)
	mockUserRepo.On("FindByID", "user-3").Return(user3, nil)

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestGetLecturerAdvisees_RepositoryError(t *testing.T) {
	service, mockLecturerRepo, mockStudentRepo, _ := setupLecturerTest()

	app := fiber.New()
	
	lecturerID := "lecturer-123"
	userID := "user-lecturer"

	app.Get("/lecturers/:id/advisees", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID: userID,
			Role:   "Admin",
		})
		return service.GetLecturerAdvisees(c)
	})

	lecturer := &model.Lecturer{
		ID:     lecturerID,
		UserID: userID,
	}

	mockLecturerRepo.On("FindByID", lecturerID).Return(lecturer, nil)
	mockStudentRepo.On("GetAll", 1000, 0).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/lecturers/"+lecturerID+"/advisees", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
	mockLecturerRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}