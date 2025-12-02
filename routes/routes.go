package routes

import (
	"project_uas/app/service"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService *service.AuthService, authMiddleware fiber.Handler) {
	auth := app.Group("/api/v1/auth")
	
	auth.Post("/login", authService.Login)
	auth.Post("/refresh", authService.Refresh)

	protected := auth.Group("/", authMiddleware)

	protected.Get("/profile", authService.Profile)
	protected.Post("/logout", authService.Logout)
}
