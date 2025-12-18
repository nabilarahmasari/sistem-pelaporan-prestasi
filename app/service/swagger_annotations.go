	package service


	// ==================== AUTH SERVICE ANNOTATIONS ======================

	// Login godoc
	// @Summary Login to the system
	// @Description Authenticate user with username/email and password
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param request body model.LoginRequest true "Login credentials"
	// @Success 200 {object} model.APIResponse{data=model.LoginResponse} "Login successful"
	// @Failure 400 {object} model.APIResponse "Invalid request body"
	// @Failure 401 {object} model.APIResponse "Invalid username or password"
	// @Router /auth/login [post]
	func (s *AuthService) LoginSwagger() {}

	// Refresh godoc
	// @Summary Refresh access token
	// @Description Get new access token using refresh token
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param request body model.RefreshTokenRequest true "Refresh token"
	// @Success 200 {object} model.APIResponse{data=model.LoginResponse} "Token refreshed"
	// @Failure 401 {object} model.APIResponse "Invalid refresh token"
	// @Failure 404 {object} model.APIResponse "User not found"
	// @Router /auth/refresh [post]
	func (s *AuthService) RefreshSwagger() {}

	// Profile godoc
	// @Summary Get current user profile
	// @Description Get profile of currently authenticated user
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} model.APIResponse{data=model.UserResponse} "User profile"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 404 {object} model.APIResponse "User not found"
	// @Router /auth/profile [get]
	func (s *AuthService) ProfileSwagger() {}

	// Logout godoc
	// @Summary Logout from system
	// @Description Logout current user (invalidate token on client side)
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} model.APIResponse "Logout successful"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Router /auth/logout [post]
	func (s *AuthService) LogoutSwagger() {}

	// ==================== USER SERVICE ANNOTATIONS ======================

	// CreateUser godoc
	// @Summary Create new user (Admin only)
	// @Description Create new user with role and profile
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body model.UserCreateRequest true "User data"
	// @Success 201 {object} model.APIResponse{data=model.UserResponse} "User created"
	// @Failure 400 {object} model.APIResponse "Invalid request body"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Failure 404 {object} model.APIResponse "Role not found"
	// @Failure 409 {object} model.APIResponse "Username/Email already exists"
	// @Failure 422 {object} model.APIResponse "Validation error"
	// @Router /users [post]
	func (s *UserService) CreateUserSwagger() {}

	// GetUsers godoc
	// @Summary Get all users (Admin only)
	// @Description Get list of all users with pagination and role filter
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param page query int false "Page number" default(1)
	// @Param page_size query int false "Page size" default(10)
	// @Param role query string false "Filter by role name" Enums(Admin, Mahasiswa, Dosen Wali)
	// @Success 200 {object} model.APIResponse{data=model.UserListResponse} "List of users"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Router /users [get]
	func (s *UserService) GetUsersSwagger() {}

	// GetUserByID godoc
	// @Summary Get user by ID (Admin only)
	// @Description Get detailed user information by ID
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID (UUID)"
	// @Success 200 {object} model.APIResponse{data=model.UserResponse} "User details"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Failure 404 {object} model.APIResponse "User not found"
	// @Router /users/{id} [get]
	func (s *UserService) GetUserByIDSwagger() {}

	// UpdateUser godoc
	// @Summary Update user (Admin only)
	// @Description Update user information
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID (UUID)"
	// @Param request body model.UserUpdateRequest true "Update data"
	// @Success 200 {object} model.APIResponse{data=model.UserResponse} "User updated"
	// @Failure 400 {object} model.APIResponse "Invalid request body"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Failure 404 {object} model.APIResponse "User not found"
	// @Failure 409 {object} model.APIResponse "Email already used"
	// @Router /users/{id} [put]
	func (s *UserService) UpdateUserSwagger() {}

	// DeleteUser godoc
	// @Summary Delete user (Admin only)
	// @Description Delete user and associated profile
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID (UUID)"
	// @Success 200 {object} model.APIResponse "User deleted"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Failure 404 {object} model.APIResponse "User not found"
	// @Router /users/{id} [delete]
	func (s *UserService) DeleteUserSwagger() {}

	// AssignRole godoc
	// @Summary Assign role to user (Admin only)
	// @Description Change user's role
	// @Tags Users
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "User ID (UUID)"
	// @Param request body model.AssignRoleRequest true "Role name"
	// @Success 200 {object} model.APIResponse{data=model.UserResponse} "Role assigned"
	// @Failure 400 {object} model.APIResponse "Invalid request body"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Failure 404 {object} model.APIResponse "User or role not found"
	// @Router /users/{id}/role [put]
	func (s *UserService) AssignRoleSwagger() {}

	// ==================== STUDENT SERVICE ANNOTATIONS ======================

	// GetAllStudents godoc
	// @Summary Get all students
	// @Description Get list of all students with pagination
	// @Tags Students
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param page query int false "Page number" default(1)
	// @Param page_size query int false "Page size" default(10)
	// @Success 200 {object} model.APIResponse "List of students with user info"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Router /students [get]
	func (s *StudentService) GetAllStudentsSwagger() {}

	// GetStudentByID godoc
	// @Summary Get student by ID
	// @Description Get detailed student information including advisor
	// @Tags Students
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Student ID (UUID)"
	// @Success 200 {object} model.APIResponse "Student details"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 404 {object} model.APIResponse "Student not found"
	// @Router /students/{id} [get]
	func (s *StudentService) GetStudentByIDSwagger() {}

	// GetStudentAchievements godoc
	// @Summary Get student achievements
	// @Description Get all achievements of a student (with authorization check)
	// @Tags Students
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Student ID (UUID)"
	// @Param page query int false "Page number" default(1)
	// @Param page_size query int false "Page size" default(10)
	// @Param status query string false "Filter by status" Enums(draft, submitted, verified, rejected)
	// @Success 200 {object} model.APIResponse "List of achievements"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized to view this student"
	// @Failure 404 {object} model.APIResponse "Student not found"
	// @Router /students/{id}/achievements [get]
	func (s *StudentService) GetStudentAchievementsSwagger() {}

	// SetAdvisor godoc
	// @Summary Set student advisor (Admin only)
	// @Description Assign or change student's advisor (lecturer)
	// @Tags Students
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Student ID (UUID)"
	// @Param request body model.SetAdvisorRequest true "Advisor ID"
	// @Success 200 {object} model.APIResponse "Advisor set successfully"
	// @Failure 400 {object} model.APIResponse "Invalid request body"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Admin only"
	// @Failure 404 {object} model.APIResponse "Student or advisor not found"
	// @Failure 422 {object} model.APIResponse "Validation error"
	// @Router /students/{id}/advisor [put]
	func (s *StudentService) SetAdvisorSwagger() {}

	// ==================== LECTURER SERVICE ANNOTATIONS ======================

	// GetAllLecturers godoc
	// @Summary Get all lecturers
	// @Description Get list of all lecturers with pagination
	// @Tags Lecturers
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param page query int false "Page number" default(1)
	// @Param page_size query int false "Page size" default(10)
	// @Success 200 {object} model.APIResponse "List of lecturers"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Router /lecturers [get]
	func (s *LecturerService) GetAllLecturersSwagger() {}

	// GetLecturerAdvisees godoc
	// @Summary Get lecturer's advisees
	// @Description Get all students advised by this lecturer
	// @Tags Lecturers
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Lecturer ID (UUID)"
	// @Success 200 {object} model.APIResponse "List of advisees"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Can only view own advisees"
	// @Failure 404 {object} model.APIResponse "Lecturer not found"
	// @Router /lecturers/{id}/advisees [get]
	func (s *LecturerService) GetLecturerAdviseesSwagger() {}

	// ==================== ACHIEVEMENT SERVICE ANNOTATIONS ======================

	// CreateAchievement godoc
	// @Summary Create achievement (Mahasiswa only)
	// @Description Create new achievement draft
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param request body model.AchievementCreateRequest true "Achievement data"
	// @Success 201 {object} model.APIResponse{data=model.AchievementResponse} "Achievement created"
	// @Failure 400 {object} model.APIResponse "Invalid request body"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Mahasiswa only"
	// @Failure 404 {object} model.APIResponse "Student profile not found"
	// @Failure 422 {object} model.APIResponse "Validation error"
	// @Router /achievements [post]
	func (s *AchievementService) CreateAchievementSwagger() {}

	// GetAchievements godoc
	// @Summary Get achievements (filtered by role)
	// @Description Get achievements list (Mahasiswa: own, Dosen: advisees, Admin: all)
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param page query int false "Page number" default(1)
	// @Param page_size query int false "Page size" default(10)
	// @Param status query string false "Filter by status" Enums(draft, submitted, verified, rejected)
	// @Success 200 {object} model.APIResponse{data=model.AchievementListResponse} "List of achievements"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden"
	// @Router /achievements [get]
	func (s *AchievementService) GetAchievementsSwagger() {}

	// GetAchievementByID godoc
	// @Summary Get achievement by ID
	// @Description Get detailed achievement information
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Success 200 {object} model.APIResponse{data=model.AchievementResponse} "Achievement details"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id} [get]
	func (s *AchievementService) GetAchievementByIDSwagger() {}

	// UpdateAchievement godoc
	// @Summary Update achievement (Mahasiswa only, draft status)
	// @Description Update achievement data (only if status is draft)
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Param request body model.AchievementUpdateRequest true "Update data"
	// @Success 200 {object} model.APIResponse{data=model.AchievementResponse} "Achievement updated"
	// @Failure 400 {object} model.APIResponse "Can only update draft achievements"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id} [put]
	func (s *AchievementService) UpdateAchievementSwagger() {}

	// DeleteAchievement godoc
	// @Summary Delete achievement (Mahasiswa only, draft status)
	// @Description Soft delete achievement (only if status is draft)
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Success 200 {object} model.APIResponse "Achievement deleted"
	// @Failure 400 {object} model.APIResponse "Can only delete draft achievements"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id} [delete]
	func (s *AchievementService) DeleteAchievementSwagger() {}

	// SubmitForVerification godoc
	// @Summary Submit achievement for verification (Mahasiswa only)
	// @Description Submit draft achievement to advisor for verification
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Success 200 {object} model.APIResponse "Achievement submitted"
	// @Failure 400 {object} model.APIResponse "Achievement already submitted"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id}/submit [post]
	func (s *AchievementService) SubmitForVerificationSwagger() {}

	// VerifyAchievement godoc
	// @Summary Verify achievement (Dosen Wali only)
	// @Description Approve submitted achievement (only for your advisees)
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Success 200 {object} model.APIResponse "Achievement verified"
	// @Failure 400 {object} model.APIResponse "Achievement must be in submitted status"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not advisor of this student"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id}/verify [post]
	func (s *AchievementService) VerifyAchievementSwagger() {}

	// RejectAchievement godoc
	// @Summary Reject achievement (Dosen Wali only)
	// @Description Reject submitted achievement with notes
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Param request body model.RejectAchievementRequest true "Rejection note"
	// @Success 200 {object} model.APIResponse "Achievement rejected"
	// @Failure 400 {object} model.APIResponse "Achievement must be in submitted status"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not advisor of this student"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Failure 422 {object} model.APIResponse "Validation error - rejection note required"
	// @Router /achievements/{id}/reject [post]
	func (s *AchievementService) RejectAchievementSwagger() {}

	// UploadAttachment godoc
	// @Summary Upload attachment file (Mahasiswa only)
	// @Description Upload file attachment to achievement (PDF, JPG, PNG max 5MB)
	// @Tags Achievements
	// @Accept multipart/form-data
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Param file formData file true "File to upload (PDF, JPG, PNG, max 5MB)"
	// @Success 201 {object} model.APIResponse{data=model.Attachment} "Attachment uploaded"
	// @Failure 400 {object} model.APIResponse "Invalid file or file too large"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not your achievement"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id}/attachments [post]
	func (s *AchievementService) UploadAttachmentSwagger() {}

	// GetAchievementHistory godoc
	// @Summary Get achievement status history
	// @Description Get timeline of achievement status changes
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Achievement Reference ID (UUID)"
	// @Success 200 {object} model.APIResponse "Achievement history"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden"
	// @Failure 404 {object} model.APIResponse "Achievement not found"
	// @Router /achievements/{id}/history [get]
	func (s *AchievementService) GetAchievementHistorySwagger() {}

	// ==================== REPORT SERVICE ANNOTATIONS ======================

	// GetStatistics godoc
	// @Summary Get achievement statistics
	// @Description Get statistics based on role (Mahasiswa: own, Dosen: advisees, Admin: all)
	// @Tags Reports
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Success 200 {object} model.APIResponse "Achievement statistics"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden"
	// @Failure 404 {object} model.APIResponse "Profile not found"
	// @Router /reports/statistics [get]
	func (s *ReportService) GetStatisticsSwagger() {}

	// GetStudentReport godoc
	// @Summary Get student achievement report
	// @Description Get comprehensive achievement report for a student
	// @Tags Reports
	// @Accept json
	// @Produce json
	// @Security BearerAuth
	// @Param id path string true "Student ID (UUID)"
	// @Success 200 {object} model.APIResponse "Student report"
	// @Failure 401 {object} model.APIResponse "Unauthorized"
	// @Failure 403 {object} model.APIResponse "Forbidden - Not authorized for this student"
	// @Failure 404 {object} model.APIResponse "Student not found"
	// @Router /reports/student/{id} [get]
	func (s *ReportService) GetStudentReportSwagger() {}