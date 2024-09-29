package middleware

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterCustomTypeFunc(validateNullString, shared.NullString{})
}

// Function to validate custom NullString type
func validateNullString(field reflect.Value) interface{} {
	if nullString, ok := field.Interface().(shared.NullString); ok {
		if nullString.Valid {
			return nullString.String
		}
	}
	return nil
}

var DTORegistry = make(map[string]map[string]interface{})

// Register DTO for routes and methods
func RegisterDTO(path, method string, dtoType interface{}) {
	if _, ok := DTORegistry[path]; !ok {
		DTORegistry[path] = make(map[string]interface{})
	}
	DTORegistry[path][method] = reflect.TypeOf(dtoType)
}

// Middleware for DTO validation
func ValidateDTO() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fullPath := c.Path()
		method := c.Method()

		// Remove "/api/v1" prefix from the path
		trimmedPath := strings.TrimPrefix(fullPath, "/api/v1")

		// Find corresponding DTO
		dtoType, found := findDTO(trimmedPath, method)
		if !found {
			return c.Next()
		}

		// Create a new instance of DTO
		dtoValue := reflect.New(dtoType).Interface()

		// Parse parameters for GET/DELETE requests and request body for other methods
		if method == "GET" || method == "DELETE" {
			if err := c.QueryParser(dtoValue); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid query parameters: " + err.Error(),
				})
			}
		} else {
			if err := c.BodyParser(dtoValue); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid data format: " + err.Error(),
				})
			}
		}

		// Check DTO for validation errors
		if err := validate.Struct(dtoValue); err != nil {
			return handleValidationError(c, err.(validator.ValidationErrors))
		}

		// Save DTO in context for further use
		c.Locals("dto", dtoValue)
		return c.Next()
	}
}

// Find DTO for specified route and method
func findDTO(fullPath, method string) (reflect.Type, bool) {
	for registeredPath, methodMap := range DTORegistry {
		if matchPath(fullPath, registeredPath) {
			if dtoType, ok := methodMap[method]; ok {
				return dtoType.(reflect.Type), true
			}
		}
	}
	return nil, false
}

// Function to match paths considering parameters
func matchPath(actualPath, registeredPath string) bool {
	actualParts := strings.Split(actualPath, "/")
	registeredParts := strings.Split(registeredPath, "/")

	if len(actualParts) != len(registeredParts) {
		return false
	}

	for i := range actualParts {
		if registeredParts[i] != actualParts[i] && !strings.HasPrefix(registeredParts[i], ":") {
			return false
		}
	}

	return true
}

// Handle validation errors
func handleValidationError(c *fiber.Ctx, errors validator.ValidationErrors) error {
	var errorMessages []string
	for _, err := range errors {
		switch err.Field() {
		case "Password":
			errorMessages = append(errorMessages, getPasswordErrorMessage(err))
		default:
			errorMessages = append(errorMessages, getErrorMessage(err))
		}
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"errors": errorMessages,
	})
}

// Returns error message for Password field
func getPasswordErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Password is required"
	case "min":
		return "Password must contain at least " + err.Param() + " characters"
	case "max":
		return "Password must contain no more than " + err.Param() + " characters"
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
	return "Password does not meet requirements: " + err.Tag()
}

// Returns general validation error message
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return "Minimum length: " + err.Param()
	case "max":
		return "Maximum length: " + err.Param()
	case "alphanum":
		return "Only letters and numbers are allowed"
	}
	return "Does not meet the rule: " + err.Tag()
}
