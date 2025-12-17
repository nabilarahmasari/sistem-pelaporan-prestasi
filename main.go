package main

import (
	"log"
	"project_uas/app/repository"
	"project_uas/routes"
	"project_uas/app/service"
	"project_uas/config"
	"project_uas/database"
	"project_uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	
	// ⭐ WAJIB: Import docs yang akan di-generate
	_ "project_uas/docs"
	
	// ⭐ WAJIB: Import fiber-swagger
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// ==================== ANOTASI GLOBAL SWAGGER ====================
// Anotasi ini mendefinisikan informasi umum API

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description REST API untuk Sistem Pelaporan Prestasi Mahasiswa dengan Role-Based Access Control (RBAC). Mendukung manajemen prestasi mahasiswa, verifikasi dosen wali, dan reporting.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support Team
// @contact.email support@unair.ac.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load config
	config.LoadEnv()
	utils.InitJWT()

	// Connect databases
	database.ConnectDatabase()
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}
	database.ConnectMongoDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(sqlDB)
	roleRepo := repository.NewRoleRepository(sqlDB)
	permRepo := repository.NewPermissionRepository(sqlDB)
	studentRepo := repository.NewStudentRepository(sqlDB)
	lecturerRepo := repository.NewLecturerRepository(sqlDB)
	achievementRepo := repository.NewAchievementRepository(sqlDB, database.MongoDB)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, permRepo)
	userService := service.NewUserService(userRepo, roleRepo, permRepo, studentRepo, lecturerRepo)
	studentService := service.NewStudentService(studentRepo, lecturerRepo, userRepo, achievementRepo)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo, userRepo)
	achievementService := service.NewAchievementService(achievementRepo, studentRepo, lecturerRepo, userRepo)
	reportService := service.NewReportService(achievementRepo, studentRepo, lecturerRepo, userRepo)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Serve static files
	app.Static("/uploads", "./uploads")

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
			"docs":    "/swagger/index.html",
		})
	})

	// ⭐ SWAGGER UI ROUTE
	// Akses: http://localhost:3000/swagger/index.html
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register API routes
	routes.AuthRoutes(app, authService)
	routes.UserRoutes(app, userService)
	routes.StudentRoutes(app, studentService)
	routes.LecturerRoutes(app, lecturerService)
	routes.AchievementRoutes(app, achievementService)
	routes.ReportRoutes(app, reportService)

	// Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server running on port %s", port)
	log.Printf("📚 Swagger docs: http://localhost:%s/swagger/index.html", port)
	log.Fatal(app.Listen(":" + port))
}