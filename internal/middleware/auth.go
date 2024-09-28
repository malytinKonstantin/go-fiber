package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/auth"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if SkipAuthMiddleware(c) {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("user_id", claims.UserID)
		return c.Next()
	}
}
