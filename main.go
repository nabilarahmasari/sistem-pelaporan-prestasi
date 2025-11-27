package main

import (
	"log"
	"os"

	"pelaporan-prestasi/config"
	"pelaporan-prestasi/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.Connect() // connect database

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	log.Println("Server running on port", port)
	app.Listen(":" + port)
}