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

	c.Locals("user", claims)
	return c.Next()
}
