package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/db"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) SetupRoutes(app *fiber.App) {
	app.Get("/users", c.ListUsers)
	app.Get("/users/:id", c.GetUser)
	app.Get("/users/username/:username", c.GetUserByUsername)
	app.Post("/users", c.CreateUser)
	app.Put("/users/:id", c.UpdateUser)
	app.Delete("/users/:id", c.DeleteUser)
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

func (c *UserController) GetUserByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	user, err := c.service.GetUserByUsername(ctx.Context(), username)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(user)
}

func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	limit, _ := strconv.ParseInt(ctx.Query("limit", "10"), 10, 32)
	offset, _ := strconv.ParseInt(ctx.Query("offset", "0"), 10, 32)

	users, err := c.service.ListUsers(ctx.Context(), int32(limit), int32(offset))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	return ctx.JSON(users)
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var input db.CreateUserParams

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := c.service.CreateUser(ctx.Context(), input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input db.UpdateUserParams
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	input.ID = int32(id)

	user, err := c.service.UpdateUser(ctx.Context(), input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return ctx.JSON(user)
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
