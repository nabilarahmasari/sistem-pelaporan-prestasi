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

func setupUserTest() (*UserService, *mocks.MockUserRepository, *mocks.MockRoleRepository, *mocks.MockPermissionRepository, *mocks.MockStudentRepository, *mocks.MockLecturerRepository) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository)
	mockPermRepo := new(mocks.MockPermissionRepository)
	mockStudentRepo := new(mocks.MockStudentRepository)
	mockLecturerRepo := new(mocks.MockLecturerRepository)

	service := NewUserService(mockUserRepo, mockRoleRepo, mockPermRepo, mockStudentRepo, mockLecturerRepo)

	return service, mockUserRepo, mockRoleRepo, mockPermRepo, mockStudentRepo, mockLecturerRepo
}

// ==================== FR-009: CREATE USER ====================

func TestCreateUser_Success_Mahasiswa(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, mockStudentRepo, mockLecturerRepo := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	roleID := "role-123"
	advisorID := "advisor-123"

	role := &model.Role{
		ID:   roleID,
		Name: "Mahasiswa",
	}

	mockUserRepo.On("FindByUsername", "mahasiswa123").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindByEmail", "mahasiswa@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("GetRoleByName", "Mahasiswa").Return(role, nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
	mockStudentRepo.On("FindByStudentID", "123456789").Return(nil, errors.New("not found"))
	mockLecturerRepo.On("FindByID", advisorID).Return(&model.Lecturer{ID: advisorID}, nil)
	mockStudentRepo.On("Create", mock.AnythingOfType("*model.Student")).Return(nil)
	mockStudentRepo.On("FindByUserID", mock.AnythingOfType("string")).Return(&model.Student{
		ID:        "student-id",
		StudentID: "123456789",
	}, nil)

	body := `{
		"username": "mahasiswa123",
		"email": "mahasiswa@test.com",
		"password": "password123",
		"full_name": "John Doe",
		"role_name": "Mahasiswa",
		"student_profile": {
			"student_id": "123456789",
			"program_study": "Teknik Informatika",
			"academic_year": "2024",
			"advisor_id": "advisor-123"
		}
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

func TestCreateUser_Success_DosenWali(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, _, mockLecturerRepo := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	roleID := "role-123"

	role := &model.Role{
		ID:   roleID,
		Name: "Dosen Wali",
	}

	mockUserRepo.On("FindByUsername", "dosen123").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindByEmail", "dosen@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("GetRoleByName", "Dosen Wali").Return(role, nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
	mockLecturerRepo.On("FindByLecturerID", "L123456").Return(nil, errors.New("not found"))
	mockLecturerRepo.On("Create", mock.AnythingOfType("*model.Lecturer")).Return(nil)
	mockLecturerRepo.On("FindByUserID", mock.AnythingOfType("string")).Return(&model.Lecturer{
		ID:         "lecturer-id",
		LecturerID: "L123456",
	}, nil)

	body := `{
		"username": "dosen123",
		"email": "dosen@test.com",
		"password": "password123",
		"full_name": "Dr. Dosen",
		"role_name": "Dosen Wali",
		"lecturer_profile": {
			"lecturer_id": "L123456",
			"department": "Informatika"
		}
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockLecturerRepo.AssertExpectations(t)
}

func TestCreateUser_Success_Admin(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	roleID := "role-admin"

	role := &model.Role{
		ID:   roleID,
		Name: "Admin",
	}

	mockUserRepo.On("FindByUsername", "admin123").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindByEmail", "admin@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("GetRoleByName", "Admin").Return(role, nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	body := `{
		"username": "admin123",
		"email": "admin@test.com",
		"password": "adminpass123",
		"full_name": "Admin User",
		"role_name": "Admin"
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestCreateUser_UsernameExists(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	existingUser := &model.User{
		ID:       "existing-id",
		Username: "mahasiswa123",
	}

	mockUserRepo.On("FindByUsername", "mahasiswa123").Return(existingUser, nil)

	body := `{
		"username": "mahasiswa123",
		"email": "new@test.com",
		"password": "password123",
		"full_name": "John Doe",
		"role_name": "Mahasiswa"
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 409, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_EmailExists(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	existingUser := &model.User{
		ID:    "existing-id",
		Email: "mahasiswa@test.com",
	}

	mockUserRepo.On("FindByUsername", "newuser").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindByEmail", "mahasiswa@test.com").Return(existingUser, nil)

	body := `{
		"username": "newuser",
		"email": "mahasiswa@test.com",
		"password": "password123",
		"full_name": "John Doe",
		"role_name": "Mahasiswa"
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 409, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestCreateUser_RoleNotFound(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	mockUserRepo.On("FindByUsername", "mahasiswa123").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindByEmail", "mahasiswa@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("GetRoleByName", "InvalidRole").Return(nil, errors.New("role not found"))

	body := `{
		"username": "mahasiswa123",
		"email": "mahasiswa@test.com",
		"password": "password123",
		"full_name": "John Doe",
		"role_name": "InvalidRole"
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestCreateUser_StudentIDExists(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, mockStudentRepo, _ := setupUserTest()

	app := fiber.New()
	app.Post("/users", service.CreateUser)

	roleID := "role-123"
	role := &model.Role{
		ID:   roleID,
		Name: "Mahasiswa",
	}

	existingStudent := &model.Student{
		ID:        "existing-student",
		StudentID: "123456789",
	}

	mockUserRepo.On("FindByUsername", "mahasiswa123").Return(nil, errors.New("not found"))
	mockUserRepo.On("FindByEmail", "mahasiswa@test.com").Return(nil, errors.New("not found"))
	mockRoleRepo.On("GetRoleByName", "Mahasiswa").Return(role, nil)
	mockUserRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
	mockStudentRepo.On("FindByStudentID", "123456789").Return(existingStudent, nil)
	mockUserRepo.On("Delete", mock.AnythingOfType("string")).Return(nil)

	body := `{
		"username": "mahasiswa123",
		"email": "mahasiswa@test.com",
		"password": "password123",
		"full_name": "John Doe",
		"role_name": "Mahasiswa",
		"student_profile": {
			"student_id": "123456789",
			"program_study": "Teknik Informatika",
			"academic_year": "2024"
		}
	}`

	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 409, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockStudentRepo.AssertExpectations(t)
}

// ==================== GET USERS ====================

func TestGetUsers_Success(t *testing.T) {
	service, mockUserRepo, _, _, mockStudentRepo, mockLecturerRepo := setupUserTest()

	app := fiber.New()
	app.Get("/users", service.GetUsers)

	users := []model.User{
		{
			ID:        "user-1",
			Username:  "user1",
			Email:     "user1@test.com",
			RoleID:    "role-1",
			CreatedAt: time.Now(),
		},
		{
			ID:        "user-2",
			Username:  "user2",
			Email:     "user2@test.com",
			RoleID:    "role-2",
			CreatedAt: time.Now(),
		},
	}

	mockUserRepo.On("GetAll", 10, 0, "").Return(users, nil)
	mockUserRepo.On("CountAll", "").Return(2, nil)
	mockUserRepo.On("GetRoleName", "role-1").Return("Mahasiswa", nil)
	mockUserRepo.On("GetRoleName", "role-2").Return("Admin", nil)
	
	// ← TAMBAH MOCK INI untuk buildUserResponse
	mockStudentRepo.On("FindByUserID", "user-1").Return(nil, errors.New("not found"))
	mockLecturerRepo.On("FindByUserID", "user-1").Return(nil, errors.New("not found"))
	mockStudentRepo.On("FindByUserID", "user-2").Return(nil, errors.New("not found"))
	mockLecturerRepo.On("FindByUserID", "user-2").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/users", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUsers_WithRoleFilter(t *testing.T) {
	service, mockUserRepo, _, _, mockStudentRepo, mockLecturerRepo := setupUserTest()

	app := fiber.New()
	app.Get("/users", service.GetUsers)

	users := []model.User{
		{
			ID:        "user-1",
			Username:  "mahasiswa1",
			Email:     "mhs1@test.com",
			RoleID:    "role-mhs",
			CreatedAt: time.Now(),
		},
	}

	mockUserRepo.On("GetAll", 10, 0, "Mahasiswa").Return(users, nil)
	mockUserRepo.On("CountAll", "Mahasiswa").Return(1, nil)
	mockUserRepo.On("GetRoleName", "role-mhs").Return("Mahasiswa", nil)
	
	// ← TAMBAH INI
	mockStudentRepo.On("FindByUserID", "user-1").Return(nil, errors.New("not found"))
	mockLecturerRepo.On("FindByUserID", "user-1").Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/users?role=Mahasiswa", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}
func TestGetUsers_WithPagination(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Get("/users", service.GetUsers)

	users := []model.User{}

	mockUserRepo.On("GetAll", 20, 20, "").Return(users, nil)
	mockUserRepo.On("CountAll", "").Return(50, nil)

	req := httptest.NewRequest("GET", "/users?page=2&page_size=20", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

// ==================== GET USER BY ID ====================

func TestGetUserByID_Success(t *testing.T) {
	service, mockUserRepo, _, _, mockStudentRepo, mockLecturerRepo := setupUserTest()

	app := fiber.New()
	app.Get("/users/:id", service.GetUserByID)

	userID := "user-123"
	user := &model.User{
		ID:        userID,
		Username:  "mahasiswa123",
		Email:     "mahasiswa@test.com",
		FullName:  "John Doe",
		RoleID:    "role-123",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockUserRepo.On("GetRoleName", "role-123").Return("Mahasiswa", nil)
	
	// ← TAMBAH INI
	mockStudentRepo.On("FindByUserID", userID).Return(nil, errors.New("not found"))
	mockLecturerRepo.On("FindByUserID", userID).Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/users/"+userID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Get("/users/:id", service.GetUserByID)

	mockUserRepo.On("FindByID", "invalid-id").Return(nil, errors.New("user not found"))

	req := httptest.NewRequest("GET", "/users/invalid-id", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

// ==================== UPDATE USER ====================

func TestUpdateUser_Success(t *testing.T) {
	service, mockUserRepo, _, _, mockStudentRepo, mockLecturerRepo := setupUserTest()

	app := fiber.New()
	app.Put("/users/:id", service.UpdateUser)

	userID := "user-123"
	user := &model.User{
		ID:        userID,
		Username:  "mahasiswa123",
		Email:     "old@test.com",
		FullName:  "Old Name",
		RoleID:    "role-123",
		CreatedAt: time.Now(),
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockUserRepo.On("FindByEmail", "new@test.com").Return(nil, errors.New("not found"))
	mockUserRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)
	mockUserRepo.On("GetRoleName", "role-123").Return("Mahasiswa", nil)
	
	// ← TAMBAH INI
	mockStudentRepo.On("FindByUserID", userID).Return(nil, errors.New("not found"))
	mockLecturerRepo.On("FindByUserID", userID).Return(nil, errors.New("not found"))

	body := `{
		"email": "new@test.com",
		"full_name": "New Name"
	}`

	req := httptest.NewRequest("PUT", "/users/"+userID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUpdateUser_EmailConflict(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Put("/users/:id", service.UpdateUser)

	userID := "user-123"
	user := &model.User{
		ID:       userID,
		Username: "mahasiswa123",
		Email:    "old@test.com",
	}

	existingUser := &model.User{
		ID:    "other-user",
		Email: "existing@test.com",
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockUserRepo.On("FindByEmail", "existing@test.com").Return(existingUser, nil)

	body := `{"email": "existing@test.com"}`

	req := httptest.NewRequest("PUT", "/users/"+userID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 409, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

// ==================== DELETE USER ====================

func TestDeleteUser_Success(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Delete("/users/:id", service.DeleteUser)

	userID := "user-123"
	user := &model.User{
		ID:       userID,
		Username: "mahasiswa123",
		RoleID:   "role-123",
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockUserRepo.On("GetRoleName", "role-123").Return("Admin", nil)
	mockUserRepo.On("Delete", userID).Return(nil)

	req := httptest.NewRequest("DELETE", "/users/"+userID, nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	service, mockUserRepo, _, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Delete("/users/:id", service.DeleteUser)

	mockUserRepo.On("FindByID", "invalid-id").Return(nil, errors.New("user not found"))

	req := httptest.NewRequest("DELETE", "/users/invalid-id", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

// ==================== ASSIGN ROLE ====================

func TestAssignRole_Success(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Put("/users/:id/role", service.AssignRole)

	userID := "user-123"
	newRoleID := "role-456"

	user := &model.User{
		ID:       userID,
		Username: "mahasiswa123",
		RoleID:   "old-role",
	}

	newRole := &model.Role{
		ID:   newRoleID,
		Name: "Admin",
	}

	updatedUser := &model.User{
		ID:        userID,
		Username:  "mahasiswa123",
		RoleID:    newRoleID,
		CreatedAt: time.Now(),
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil).Once()
	mockRoleRepo.On("GetRoleByName", "Admin").Return(newRole, nil)
	mockUserRepo.On("UpdateRole", userID, newRoleID).Return(nil)
	mockUserRepo.On("FindByID", userID).Return(updatedUser, nil).Once()

	body := `{"role_name": "Admin"}`

	req := httptest.NewRequest("PUT", "/users/"+userID+"/role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestAssignRole_InvalidRole(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, _, _, _ := setupUserTest()

	app := fiber.New()
	app.Put("/users/:id/role", service.AssignRole)

	userID := "user-123"
	user := &model.User{
		ID:       userID,
		Username: "mahasiswa123",
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockRoleRepo.On("GetRoleByName", "InvalidRole").Return(nil, errors.New("role not found"))

	body := `{"role_name": "InvalidRole"}`

	req := httptest.NewRequest("PUT", "/users/"+userID+"/role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}