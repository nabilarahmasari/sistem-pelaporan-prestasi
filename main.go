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
)

func main() {
	// Load config
	config.LoadEnv()

	// Initialize JWT
	utils.InitJWT()

	// Connect PostgreSQL database
	database.ConnectDatabase()
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Failed to get database connection:", err)
	}

	// Connect MongoDB database
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
	
	// StudentService dengan achievementRepo
	studentService := service.NewStudentService(studentRepo, lecturerRepo, userRepo, achievementRepo)
	
	// LecturerService (NEW)
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

	// Serve static files dari folder uploads
	app.Static("/uploads", "./uploads")

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	// Register routes
	routes.AuthRoutes(app, authService)
	routes.UserRoutes(app, userService)
	routes.StudentRoutes(app, studentService)
	routes.LecturerRoutes(app, lecturerService) // NEW
	routes.AchievementRoutes(app, achievementService)
	routes.ReportRoutes(app, reportService) // NEW

	// Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}