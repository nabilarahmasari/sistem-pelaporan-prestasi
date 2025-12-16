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
	"project_uas/utils"
)

// ==================== HELPER FUNCTIONS ====================

func setupAuthTest() (*AuthService, *mocks.MockUserRepository, *mocks.MockRoleRepository, *mocks.MockPermissionRepository) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockRoleRepo := new(mocks.MockRoleRepository)
	mockPermRepo := new(mocks.MockPermissionRepository)

	service := NewAuthService(mockUserRepo, mockRoleRepo, mockPermRepo)

	return service, mockUserRepo, mockRoleRepo, mockPermRepo
}

// ==================== FR-001: LOGIN ====================

func TestLogin_Success(t *testing.T) {
	service, mockUserRepo, mockRoleRepo, mockPermRepo := setupAuthTest()

	app := fiber.New()
	app.Post("/login", service.Login)

	userID := "user-123"
	roleID := "role-123"
	hashedPassword, _ := utils.HashPassword("password123")

	user := &model.User{
		ID:           userID,
		Username:     "mahasiswa123",
		Email:        "mahasiswa@test.com",
		PasswordHash: hashedPassword,
		FullName:     "John Doe",
		RoleID:       roleID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	role := &model.Role{
		ID:   roleID,
		Name: "Mahasiswa",
	}

	permissions := []string{"achievement:create", "achievement:read"}

	// Mock expectations
	mockUserRepo.On("FindByUsername", "mahasiswa123").Return(user, nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockPermRepo.On("GetPermissionsByRoleID", roleID).Return(permissions, nil)

	body := `{"username": "mahasiswa123", "password": "password123"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockPermRepo.AssertExpectations(t)
}

func TestLogin_InvalidUsername(t *testing.T) {
	service, mockUserRepo, _, _ := setupAuthTest()

	app := fiber.New()
	app.Post("/login", service.Login)

	mockUserRepo.On("FindByUsername", "invaliduser").Return(nil, errors.New("user not found"))

	body := `{"username": "invaliduser", "password": "password123"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	service, mockUserRepo, _, _ := setupAuthTest()

	app := fiber.New()
	app.Post("/login", service.Login)

	userID := "user-123"
	hashedPassword, _ := utils.HashPassword("correctpassword")

	user := &model.User{
		ID:           userID,
		Username:     "mahasiswa123",
		PasswordHash: hashedPassword,
	}

	mockUserRepo.On("FindByUsername", "mahasiswa123").Return(user, nil)

	body := `{"username": "mahasiswa123", "password": "wrongpassword"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	service, _, _, _ := setupAuthTest()

	app := fiber.New()
	app.Post("/login", service.Login)

	body := `{"invalid": "json"`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

// ==================== REFRESH TOKEN ====================

func TestRefreshToken_Success(t *testing.T) {
	t.Skip("Skipped: Refresh token validation requires JWT secret from config/environment. This is an integration test that should be run with proper environment setup.")
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	service, _, _, _ := setupAuthTest()

	app := fiber.New()
	app.Post("/refresh", service.Refresh)

	body := `{"refresh_token": "invalid.token.here"}`
	req := httptest.NewRequest("POST", "/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	t.Skip("Skipped: Refresh token validation requires JWT secret from config/environment. This is an integration test that should be run with proper environment setup.")
}

// ==================== PROFILE ====================

func TestProfile_Success(t *testing.T) {
	service, mockUserRepo, _, _ := setupAuthTest()

	app := fiber.New()
	
	userID := "user-123"

	app.Get("/profile", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID:   userID,
			Username: "mahasiswa123",
			Role:     "Mahasiswa",
		})
		return service.Profile(c)
	})

	user := &model.User{
		ID:        userID,
		Username:  "mahasiswa123",
		Email:     "mahasiswa@test.com",
		FullName:  "John Doe",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	mockUserRepo.On("FindByID", userID).Return(user, nil)

	req := httptest.NewRequest("GET", "/profile", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

// ==================== LOGOUT ====================

func TestLogout_Success(t *testing.T) {
	service, _, _, _ := setupAuthTest()

	app := fiber.New()
	
	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Locals("user", &model.JWTClaims{
			UserID:   "user-123",
			Username: "mahasiswa123",
			Role:     "Mahasiswa",
		})
		return service.Logout(c)
	})

	req := httptest.NewRequest("POST", "/logout", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestLogout_Unauthorized(t *testing.T) {
	service, _, _, _ := setupAuthTest()

	app := fiber.New()
	app.Post("/logout", service.Logout)

	req := httptest.NewRequest("POST", "/logout", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 401, resp.StatusCode)
}