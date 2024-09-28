package user

import (
	"github.com/gofiber/fiber/v2"

	"strconv"
)

func SetupRoutes(app *fiber.App, controller *UserController) {
	app.Get("/users", controller.ListUsers)
	app.Get("/users/:id", controller.GetUser)
	app.Post("/users", controller.CreateUser)
	app.Delete("/users/:id", controller.DeleteUser)
}

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user, err := c.service.GetUser(ctx.Context(), int32(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(user)
}

func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	users, err := c.service.ListUsers(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	return ctx.JSON(users)
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := c.service.CreateUser(ctx.Context(), input.Name, input.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := c.service.DeleteUser(ctx.Context(), int32(id)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
