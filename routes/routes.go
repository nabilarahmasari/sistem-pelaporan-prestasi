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

	// Auth required untuk semua endpoint
	students.Use(middleware.AuthRequired)

	students.Get("/", studentService.GetAllStudents)           // GET /api/v1/students
	students.Get("/:id", studentService.GetStudentByID)        // GET /api/v1/students/:id
	students.Put("/:id/advisor", studentService.SetAdvisor)    // PUT /api/v1/students/:id/advisor
	
	// TODO: GET /api/v1/students/:id/achievements (nanti saat implement achievements)
}

//
// ==================== LECTURER ROUTES (OPTIONAL) ======================
//

// func LecturerRoutes(app *fiber.App, lecturerService *service.LecturerService) {
// 	lecturers := app.Group("/api/v1/lecturers")
// 	lecturers.Use(middleware.AuthRequired)
//
// 	lecturers.Get("/", lecturerService.GetAllLecturers)
// 	lecturers.Get("/:id/advisees", lecturerService.GetAdvisees)
// }