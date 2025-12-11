package routes

import (
	"project_uas/middleware"
	"project_uas/app/service"

	"github.com/gofiber/fiber/v2"
)

//
// ==================== AUTH ROUTES ======================
//

func AuthRoutes(app *fiber.App, authService *service.AuthService) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.Refresh)

	// Protected routes
	protected := auth.Group("/", middleware.AuthRequired)
	protected.Get("/profile", authService.Profile)
	protected.Post("/logout", authService.Logout)
}

//
// ==================== USER ROUTES (ADMIN ONLY) ======================
//

func UserRoutes(app *fiber.App, userService *service.UserService) {
	users := app.Group("/api/v1/users")

	// Semua endpoint user butuh auth + permission "user:manage"
	users.Use(middleware.AuthRequired)
	users.Use(middleware.RequirePermission("user:manage"))

	users.Get("/", userService.GetUsers)          // GET /api/v1/users
	users.Get("/:id", userService.GetUserByID)    // GET /api/v1/users/:id
	users.Post("/", userService.CreateUser)       // POST /api/v1/users
	users.Put("/:id", userService.UpdateUser)     // PUT /api/v1/users/:id
	users.Delete("/:id", userService.DeleteUser)  // DELETE /api/v1/users/:id
	users.Put("/:id/role", userService.AssignRole) // PUT /api/v1/users/:id/role
}
//
// ==================== STUDENT ROUTES ======================
//

func StudentRoutes(app *fiber.App, studentService *service.StudentService) {
	students := app.Group("/api/v1/students")
	students.Use(middleware.AuthRequired)

	// GET /students - List all students
	students.Get("/",
		studentService.GetAllStudents,
	)

	// GET /students/:id - Get student by ID
	students.Get("/:id",
		studentService.GetStudentByID,
	)

	// GET /students/:id/achievements - Get student achievements
	students.Get("/:id/achievements",
		middleware.RequirePermission("achievement:read"),
		studentService.GetStudentAchievements,
	)

	// PUT /students/:id/advisor - Set advisor (Admin only)
	students.Put("/:id/advisor",
		middleware.RequirePermission("user:manage"),
		studentService.SetAdvisor,
	)
}

//
// ==================== LECTURER ROUTES ======================
//

func LecturerRoutes(app *fiber.App, lecturerService *service.LecturerService) {
	lecturers := app.Group("/api/v1/lecturers")
	lecturers.Use(middleware.AuthRequired)

	// GET /lecturers - List all lecturers
	lecturers.Get("/",
		lecturerService.GetAllLecturers,
	)

	// GET /lecturers/:id/advisees - Get lecturer's advisees
	lecturers.Get("/:id/advisees",
		lecturerService.GetLecturerAdvisees,
	)
}
//
// ==================== ACHIEVEMENT ROUTES ======================
//

func AchievementRoutes(app *fiber.App, achievementService *service.AchievementService) {
	achievements := app.Group("/api/v1/achievements")

	// Auth required untuk semua endpoint
	achievements.Use(middleware.AuthRequired)

	// GET /achievements - List achievements (filtered by role)
	// Mahasiswa: achievement:read (own)
	// Dosen Wali: achievement:read (advisees)
	// Admin: achievement:read (all)
	achievements.Get("/", 
		middleware.RequirePermission("achievement:read"),
		achievementService.GetAchievements,
	)

	// GET /achievements/:id - Detail achievement
	achievements.Get("/:id",
		middleware.RequirePermission("achievement:read"),
		achievementService.GetAchievementByID,
	)

	// POST /achievements - Create achievement (Mahasiswa only)
	achievements.Post("/",
		middleware.RequirePermission("achievement:create"),
		achievementService.CreateAchievement,
	)

	// PUT /achievements/:id - Update achievement (Mahasiswa only, status = draft)
	achievements.Put("/:id",
		middleware.RequirePermission("achievement:update"),
		achievementService.UpdateAchievement,
	)

	// DELETE /achievements/:id - Delete achievement (Mahasiswa only, status = draft)
	achievements.Delete("/:id",
		middleware.RequirePermission("achievement:delete"),
		achievementService.DeleteAchievement,
	)

	// POST /achievements/:id/submit - Submit for verification (Mahasiswa only)
	achievements.Post("/:id/submit",
		middleware.RequirePermission("achievement:update"),
		achievementService.SubmitForVerification,
	)

	// POST /achievements/:id/verify - Verify achievement (Dosen Wali only)
	achievements.Post("/:id/verify",
		middleware.RequirePermission("achievement:verify"),
		achievementService.VerifyAchievement,
	)

	// POST /achievements/:id/reject - Reject achievement (Dosen Wali only)
	achievements.Post("/:id/reject",
		middleware.RequirePermission("achievement:verify"),
		achievementService.RejectAchievement,
	)

	// POST /achievements/:id/attachments - Upload attachment (Mahasiswa only)
	achievements.Post("/:id/attachments",
		middleware.RequirePermission("achievement:update"),
		achievementService.UploadAttachment,
	)

	// GET /achievements/:id/history - History achievement
    achievements.Get("/:id/history",
        middleware.RequirePermission("achievement:read"),
        achievementService.GetAchievementHistory,
    )
}