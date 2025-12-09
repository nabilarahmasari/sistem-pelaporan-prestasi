package middleware

import (
	"strings"
	"project_uas/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing token"})
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}

	// ⭐ SET USER CLAIMS
	c.Locals("user", claims)
	
	// ⭐ SET PERMISSIONS - INI YANG KURANG!
	c.Locals("permissions", claims.Permissions)
	
	return c.Next()
}

func RequirePermission(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "error",
				"error":  "Akses ditolak: permission tidak tersedia",
			})
		}

		for _, perm := range permissions {
			if perm == requiredPermission {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "error",
			"error":  "Akses ditolak: permission tidak mencukupi",
		})
	}
}