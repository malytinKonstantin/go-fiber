package user

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/db"
)

// swagger:model FiberMap
type FiberMap map[string]interface{}

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

// SetupRoutes sets up the user-related routes
// @Summary Set up user routes
// @Description Set up routes for user-related operations
// @Tags users
func (c *UserController) SetupRoutes(app *fiber.App) {
	app.Get("/users", c.ListUsers)
	app.Get("/users/:id", c.GetUser)
	app.Get("/users/username/:username", c.GetUserByUsername)
	app.Post("/users", c.CreateUser)
	app.Put("/users/:id", c.UpdateUser)
	app.Delete("/users/:id", c.DeleteUser)
}

// GetUser retrieves a user by ID
// @Summary Get a user
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} FiberMap
// @Failure 404 {object} FiberMap
// @Router /users/{id} [get]
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

// GetUserByUsername retrieves a user by username
// @Summary Get a user by username
// @Description Get a user by their username
// @Tags users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} User
// @Failure 404 {object} FiberMap
// @Router /users/username/{username} [get]
func (c *UserController) GetUserByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	user, err := c.service.GetUserByUsername(ctx.Context(), username)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(user)
}

// ListUsers retrieves a list of users
// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} User
// @Failure 500 {object} FiberMap
// @Router /users [get]
func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	limit, _ := strconv.ParseInt(ctx.Query("limit", "10"), 10, 32)
	offset, _ := strconv.ParseInt(ctx.Query("offset", "0"), 10, 32)

	users, err := c.service.ListUsers(ctx.Context(), int32(limit), int32(offset))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	return ctx.JSON(users)
}

// CreateUser creates a new user
// @Summary Create a user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserParams true "User information"
// @Success 201 {object} User
// @Failure 400 {object} FiberMap
// @Failure 500 {object} FiberMap
// @Router /users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var input CreateUserParams

	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := c.service.CreateUser(ctx.Context(), input)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body db.UpdateUserParams true "Updated user information"
// @Success 200 {object} User
// @Failure 400 {object} FiberMap
// @Failure 500 {object} FiberMap
// @Router /users/{id} [put]
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

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} FiberMap
// @Failure 500 {object} FiberMap
// @Router /users/{id} [delete]
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
