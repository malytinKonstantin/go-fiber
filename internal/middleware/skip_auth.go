package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func SkipAuthMiddleware(c *fiber.Ctx) bool {
	return c.Locals("skip_auth") == true
}

func SkipAuth(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("skip_auth", true)
		return handler(c)
	}
}
