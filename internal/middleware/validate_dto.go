package middleware

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterCustomTypeFunc(validateNullString, shared.NullString{})
}

func validateNullString(field reflect.Value) interface{} {
	if nullString, ok := field.Interface().(shared.NullString); ok {
		if nullString.Valid {
			return nullString.String
		}
	}
	return nil
}

var DTORegistry = make(map[string]map[string]interface{})

func RegisterDTO(path, method string, dtoType interface{}) {
	if _, ok := DTORegistry[path]; !ok {
		DTORegistry[path] = make(map[string]interface{})
	}
	DTORegistry[path][method] = reflect.TypeOf(dtoType)
}

func ValidateDTO() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fullPath := c.Path()
		method := c.Method()

		fmt.Printf("ValidateDTO: FullPath=%s, Method=%s\n", fullPath, method)
		fmt.Printf("DTORegistry: %+v\n", DTORegistry)

		dtoType, ok := DTORegistry[fullPath][method]
		if !ok {
			fmt.Printf("No DTO registered for path=%s, method=%s\n", fullPath, method)
			return c.Next()
		}

		dtoValue := reflect.New(dtoType.(reflect.Type)).Interface()

		if err := c.BodyParser(dtoValue); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid data format: " + err.Error(),
			})
		}

		if err := validate.Struct(dtoValue); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				switch err.Field() {
				case "Password":
					errors = append(errors, getPasswordErrorMessage(err))
				default:
					errors = append(errors, fmt.Sprintf("Field '%s': %s", err.Field(), getErrorMessage(err)))
				}
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": errors,
			})
		}

		c.Locals("dto", dtoValue)
		fmt.Printf("DTO set in context: %+v\n", dtoValue)

		return c.Next()
	}
}

func getPasswordErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Password is required"
	case "min":
		return fmt.Sprintf("Password must contain at least %s characters", err.Param())
	case "max":
		return fmt.Sprintf("Password must contain no more than %s characters", err.Param())
	case "containsany":
		switch err.Param() {
		case "abcdefghijklmnopqrstuvwxyz":
			return "Password must contain at least one lowercase letter"
		case "ABCDEFGHIJKLMNOPQRSTUVWXYZ":
			return "Password must contain at least one uppercase letter"
		case "0123456789":
			return "Password must contain at least one digit"
		case "!@#$%^&*()":
			return "Password must contain at least one special character (!@#$%^&*())"
		}
	}
	return fmt.Sprintf("Password does not meet requirements: %s", err.Tag())
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return fmt.Sprintf("Minimum length: %s", err.Param())
	case "max":
		return fmt.Sprintf("Maximum length: %s", err.Param())
	case "alphanum":
		return "Only letters and numbers are allowed"
	}
	return fmt.Sprintf("Does not meet the rule: %s", err.Tag())
}
