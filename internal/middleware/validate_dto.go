package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ValidateDTO(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(dto); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid data format",
			})
		}

		if err := validate.Struct(dto); err != nil {
			errors := err.(validator.ValidationErrors)
			errorMessages := make([]string, 0)

			for _, e := range errors {
				errorMessages = append(errorMessages, formatError(e))
			}

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": errorMessages,
			})
		}

		c.Locals("dto", dto)
		return c.Next()
	}
}

func formatError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return err.Field() + " must be a valid email address"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters long"
	case "max":
		return err.Field() + " must be no more than " + err.Param() + " characters long"
	case "gte":
		return err.Field() + " must be greater than or equal to " + err.Param()
	case "lte":
		return err.Field() + " must be less than or equal to " + err.Param()
	default:
		return err.Field() + " is invalid"
	}
}
