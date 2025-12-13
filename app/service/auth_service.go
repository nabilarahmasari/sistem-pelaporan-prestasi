package service

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"project_uas/app/model"
	"project_uas/app/repository"
	"project_uas/utils"
)

type AuthService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	permRepo repository.PermissionRepository
}

func NewAuthService(
	user repository.UserRepository,
	role repository.RoleRepository,
	perm repository.PermissionRepository,
) *AuthService {
	return &AuthService{
		userRepo: user,
		roleRepo: role,
		permRepo: perm,
	}
}

//
// ==================== LOGIN ======================
//

func (s *AuthService) Login(c *fiber.Ctx) error {

	req := new(model.LoginRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// cek username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid username or password",
		})
	}

	// cek password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid username or password",
		})
	}

	// ambil role
	role, _ := s.roleRepo.GetRoleByID(user.RoleID)

	// ambil permission by role
	perms, _ := s.permRepo.GetPermissionsByRoleID(role.ID)

	// response user
	userRes := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        role.Name,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		Permissions: perms,
	}

	// generate token
	access, _ := utils.GenerateJWT(userRes)
	refresh, _ := utils.GenerateRefreshToken(user.ID)

	// Log successful login
	log.Printf("[LOGIN] User: %s (%s) | Role: %s | Time: %s",
		user.Username,
		user.ID,
		role.Name,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: model.LoginResponse{
			Token:        access,
			RefreshToken: refresh,
			User:         userRes,
		},
	})
}

//
// ==================== REFRESH TOKEN ======================
//

func (s *AuthService) Refresh(c *fiber.Ctx) error {

	req := new(model.RefreshTokenRequest)
	_ = c.BodyParser(req)

	claims := &model.RefreshTokenClaims{}

	// validasi refresh token
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid refresh token",
		})
	}

	// cari user
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "user not found",
		})
	}

	// ambil role
	role, _ := s.roleRepo.GetRoleByID(user.RoleID)

	// ambil permission
	perms, _ := s.permRepo.GetPermissionsByRoleID(role.ID)

	userRes := model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        role.Name,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
		Permissions: perms,
	}

	// generate new tokens
	access, _ := utils.GenerateJWT(userRes)
	refresh, _ := utils.GenerateRefreshToken(user.ID)

	// Log token refresh
	log.Printf("[REFRESH] User: %s (%s) | Time: %s",
		user.Username,
		user.ID,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data: model.LoginResponse{
			Token:        access,
			RefreshToken: refresh,
			User:         userRes,
		},
	})
}

//
// ==================== PROFILE ======================
//

func (s *AuthService) Profile(c *fiber.Ctx) error {

	claims := c.Locals("user").(*model.JWTClaims)

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "user not found",
		})
	}

	res := model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      claims.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   res,
	})
}

//
// ==================== LOGOUT (IMPROVED VERSION) ======================
//

func (s *AuthService) Logout(c *fiber.Ctx) error {
	// Verify user authenticated (defensive programming)
	claims, ok := c.Locals("user").(*model.JWTClaims)
	if !ok {
		return c.Status(401).JSON(model.APIResponse{
			Status: "error",
			Error:  "unauthorized",
		})
	}

	// Log logout activity (audit trail)
	log.Printf("[LOGOUT] User: %s (%s) | Role: %s | Time: %s",
		claims.Username,
		claims.UserID,
		claims.Role,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "logout successful",
	})
}