package service

import (
	"math"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"project_uas/app/model"
	"project_uas/app/repository"
	"project_uas/utils"
)

type UserService struct {
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	permRepo     repository.PermissionRepository
	studentRepo  repository.StudentRepository
	lecturerRepo repository.LecturerRepository
	validate     *validator.Validate
}

func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		permRepo:     permRepo,
		studentRepo:  studentRepo,
		lecturerRepo: lecturerRepo,
		validate:     validator.New(),
	}
}

//
// ==================== CREATE USER (POST /users) ======================
//

func (s *UserService) CreateUser(c *fiber.Ctx) error {
	req := new(model.UserCreateRequest)

	// Parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Validasi input
	if err := s.validate.Struct(req); err != nil {
		return c.Status(422).JSON(model.APIResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	// Cek username sudah ada atau belum
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return c.Status(409).JSON(model.APIResponse{
			Status: "error",
			Error:  "username already exists",
		})
	}

	// Cek email sudah ada atau belum
	existingEmail, _ := s.userRepo.FindByEmail(req.Email)
	if existingEmail != nil {
		return c.Status(409).JSON(model.APIResponse{
			Status: "error",
			Error:  "email already exists",
		})
	}

	// Get role by name
	role, err := s.roleRepo.GetRoleByName(req.RoleName)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "role not found",
		})
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to hash password",
		})
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		RoleID:       role.ID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to create user",
		})
	}

	// Create profile berdasarkan role
	if role.Name == "Mahasiswa" && req.StudentProfile != nil {
		// Validasi student profile
		if err := s.validate.Struct(req.StudentProfile); err != nil {
			// Rollback user creation (delete user)
			s.userRepo.Delete(user.ID)
			return c.Status(422).JSON(model.APIResponse{
				Status: "error",
				Error:  "invalid student profile: " + err.Error(),
			})
		}

		// Cek student_id sudah ada atau belum
		existingStudent, _ := s.studentRepo.FindByStudentID(req.StudentProfile.StudentID)
		if existingStudent != nil {
			s.userRepo.Delete(user.ID)
			return c.Status(409).JSON(model.APIResponse{
				Status: "error",
				Error:  "student_id already exists",
			})
		}

		// Jika ada advisor_id, validasi lecturer exists
		if req.StudentProfile.AdvisorID != nil {
			_, err := s.lecturerRepo.FindByID(*req.StudentProfile.AdvisorID)
			if err != nil {
				s.userRepo.Delete(user.ID)
				return c.Status(404).JSON(model.APIResponse{
					Status: "error",
					Error:  "advisor not found",
				})
			}
		}

		// Create student profile
		student := &model.Student{
			UserID:       user.ID,
			StudentID:    req.StudentProfile.StudentID,
			ProgramStudy: req.StudentProfile.ProgramStudy,
			AcademicYear: req.StudentProfile.AcademicYear,
			AdvisorID:    req.StudentProfile.AdvisorID,
		}

		if err := s.studentRepo.Create(student); err != nil {
			s.userRepo.Delete(user.ID)
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to create student profile",
			})
		}
	}

	if role.Name == "Dosen Wali" && req.LecturerProfile != nil {
		// Validasi lecturer profile
		if err := s.validate.Struct(req.LecturerProfile); err != nil {
			s.userRepo.Delete(user.ID)
			return c.Status(422).JSON(model.APIResponse{
				Status: "error",
				Error:  "invalid lecturer profile: " + err.Error(),
			})
		}

		// Cek lecturer_id sudah ada atau belum
		existingLecturer, _ := s.lecturerRepo.FindByLecturerID(req.LecturerProfile.LecturerID)
		if existingLecturer != nil {
			s.userRepo.Delete(user.ID)
			return c.Status(409).JSON(model.APIResponse{
				Status: "error",
				Error:  "lecturer_id already exists",
			})
		}

		// Create lecturer profile
		lecturer := &model.Lecturer{
			UserID:     user.ID,
			LecturerID: req.LecturerProfile.LecturerID,
			Department: req.LecturerProfile.Department,
		}

		if err := s.lecturerRepo.Create(lecturer); err != nil {
			s.userRepo.Delete(user.ID)
			return c.Status(500).JSON(model.APIResponse{
				Status: "error",
				Error:  "failed to create lecturer profile",
			})
		}
	}

	// Build response dengan profile
	userResponse := s.buildUserResponse(user, role.Name)

	return c.Status(201).JSON(model.APIResponse{
		Status:  "success",
		Message: "user created successfully",
		Data:    userResponse,
	})
}

//
// ==================== GET ALL USERS (GET /users) ======================
//

func (s *UserService) GetUsers(c *fiber.Ctx) error {
	// Parse query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	roleName := c.Query("role", "") // filter by role (optional)

	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get users from repository
	users, err := s.userRepo.GetAll(pageSize, offset, roleName)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to fetch users",
		})
	}

	// Get total count
	total, err := s.userRepo.CountAll(roleName)
	if err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to count users",
		})
	}

	// Build response
	var userResponses []model.UserResponse
	for _, user := range users {
		roleName, _ := s.userRepo.GetRoleName(user.RoleID)
		userResp := s.buildUserResponse(&user, roleName)
		userResponses = append(userResponses, *userResp)
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := model.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   response,
	})
}

//
// ==================== GET USER BY ID (GET /users/:id) ======================
//

func (s *UserService) GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "user not found",
		})
	}

	roleName, _ := s.userRepo.GetRoleName(user.RoleID)
	userResponse := s.buildUserResponse(user, roleName)

	return c.JSON(model.APIResponse{
		Status: "success",
		Data:   userResponse,
	})
}

//
// ==================== UPDATE USER (PUT /users/:id) ======================
//

func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Cari user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "user not found",
		})
	}

	// Parse request
	req := new(model.UserUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Update fields (hanya yang diisi)
	if req.Email != "" {
		// Cek email conflict
		existingEmail, _ := s.userRepo.FindByEmail(req.Email)
		if existingEmail != nil && existingEmail.ID != userID {
			return c.Status(409).JSON(model.APIResponse{
				Status: "error",
				Error:  "email already used by another user",
			})
		}
		user.Email = req.Email
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update ke database
	if err := s.userRepo.Update(user); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to update user",
		})
	}

	roleName, _ := s.userRepo.GetRoleName(user.RoleID)
	userResponse := s.buildUserResponse(user, roleName)

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "user updated successfully",
		Data:    userResponse,
	})
}

//
// ==================== DELETE USER (DELETE /users/:id) ======================
//

func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Cari user dulu
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "user not found",
		})
	}

	// Hapus profile dulu (jika ada) - CASCADE DELETE seharusnya handle ini
	roleName, _ := s.userRepo.GetRoleName(user.RoleID)

	if roleName == "Mahasiswa" {
		student, _ := s.studentRepo.FindByUserID(userID)
		if student != nil {
			s.studentRepo.Delete(student.ID)
		}
	}

	if roleName == "Dosen Wali" {
		lecturer, _ := s.lecturerRepo.FindByUserID(userID)
		if lecturer != nil {
			s.lecturerRepo.Delete(lecturer.ID)
		}
	}

	// Hapus user
	if err := s.userRepo.Delete(userID); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to delete user",
		})
	}

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "user deleted successfully",
	})
}

//
// ==================== ASSIGN ROLE (PUT /users/:id/role) ======================
//

func (s *UserService) AssignRole(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Cari user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "user not found",
		})
	}

	// Parse request
	req := new(model.AssignRoleRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(model.APIResponse{
			Status: "error",
			Error:  "invalid request body",
		})
	}

	// Get role by name
	role, err := s.roleRepo.GetRoleByName(req.RoleName)
	if err != nil {
		return c.Status(404).JSON(model.APIResponse{
			Status: "error",
			Error:  "role not found",
		})
	}

	// Update role
	if err := s.userRepo.UpdateRole(userID, role.ID); err != nil {
		return c.Status(500).JSON(model.APIResponse{
			Status: "error",
			Error:  "failed to assign role",
		})
	}

	// Refresh user data
	user, _ = s.userRepo.FindByID(userID)
	userResponse := s.buildUserResponse(user, role.Name)

	return c.JSON(model.APIResponse{
		Status:  "success",
		Message: "role assigned successfully",
		Data:    userResponse,
	})
}

//
// ==================== HELPER: BUILD USER RESPONSE ======================
//

func (s *UserService) buildUserResponse(user *model.User, roleName string) *model.UserResponse {
	response := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      roleName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// Load profile jika ada
	if roleName == "Mahasiswa" {
		student, err := s.studentRepo.FindByUserID(user.ID)
		if err == nil {
			response.StudentProfile = &model.StudentResponse{
				ID:           student.ID,
				UserID:       student.UserID,
				StudentID:    student.StudentID,
				ProgramStudy: student.ProgramStudy,
				AcademicYear: student.AcademicYear,
				AdvisorID:    student.AdvisorID,
				CreatedAt:    student.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}
	}

	if roleName == "Dosen Wali" {
		lecturer, err := s.lecturerRepo.FindByUserID(user.ID)
		if err == nil {
			response.LecturerProfile = &model.LecturerResponse{
				ID:         lecturer.ID,
				UserID:     lecturer.UserID,
				LecturerID: lecturer.LecturerID,
				Department: lecturer.Department,
				CreatedAt:  lecturer.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}
	}

	return response
}
