package main

import (
	"database/sql"
	"log"

	"project_uas/config"
	"project_uas/database"
	"project_uas/routes"
	"project_uas/app/repository"
	"project_uas/app/service"
	"project_uas/middleware"

	_ "github.com/lib/pq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config.LoadEnv()

	// Connect database (sql.DB)
	db, err := sql.Open("postgres", config.AppConfig.DBUrl)
	if err != nil {
		log.Fatal("Failed connect:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed ping database:", err)
	}
	log.Println("✅ Database connected!")

	// Run migration
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Run seeder
	if err := database.RunSeeders(db); err != nil {
		log.Fatal("Seeder failed:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"status": "error",
				"error":  err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// ⭐ PENTING: Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	// ⭐ PENTING: Initialize Services
	authService := service.NewAuthService(userRepo, roleRepo, permRepo)

	// ⭐ PENTING: Initialize Middleware
	authMiddleware := middleware.AuthRequired

	// ⭐ PENTING: Setup Routes (INI YANG KURANG!)
	routes.AuthRoutes(app, authService, authMiddleware)

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	// Start server
	log.Println("🚀 Server starting on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}