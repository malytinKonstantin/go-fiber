package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/middleware"
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
func (c *UserController) SetupRoutes(router fiber.Router) {
	// public routes
	router.Post("/signin", middleware.SkipAuth(c.SignIn))
	middleware.RegisterDTO("/api/v1/signin", "POST", SignInDto{})
	router.Post("/signup", middleware.SkipAuth(c.CreateUser))

	// protected routes
	router.Get("/users", c.ListUsers)
	router.Get("/users/:id", c.GetUser)
	router.Get("/users/username/:username", c.GetUserByUsername)
	middleware.RegisterDTO("/api/v1/users", "POST", CreateUserDto{})
	router.Post("/users", c.CreateUser)
	middleware.RegisterDTO("/api/v1/users", "PATCH", UpdateUserDto{})
	router.Patch("/users/:id", c.UpdateUser)
	router.Delete("/users/:id", c.DeleteUser)
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
// @Router /api/v1/users/{id} [get]
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
// @Router /api/v1/users/username/{username} [get]
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
// @Description Get a list of users with pagination and search options
// @Tags users
// @Accept json
// @Produce json
// @Param username query string false "Username"
// @Param email query string false "Email"
// @Param full_name query string false "Full Name"
// @Param bio query string false "Bio"
// @Param created_from query string false "Created From (RFC3339 format)"
// @Param created_to query string false "Created To (RFC3339 format)"
// @Param sort_by query string false "Sort By"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} User
// @Failure 500 {object} FiberMap
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	query := new(ListUsersQuery)

	if err := ctx.QueryParser(query); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	params := SearchUsersParams{
		Limit:  100,
		Offset: 0,
	}

	if query.Username != "" {
		params.Username = query.Username
	}

	if query.Email != "" {
		params.Email = query.Email
	}

	if query.FullName != "" {
		params.FullName = query.FullName
	}

	if query.Bio != "" {
		params.Bio = query.Bio
	}

	if query.CreatedFrom != "" {
		params.CreatedFrom = query.CreatedFrom
	}

	if query.CreatedTo != "" {
		params.CreatedTo = query.CreatedTo
	}

	if query.SortBy != "" {
		params.SortBy = query.SortBy
	}

	if query.Limit > 0 {
		params.Limit = query.Limit
	}

	if query.Offset >= 0 {
		params.Offset = query.Offset
	}

	users, err := c.service.SearchUsers(ctx.Context(), params)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	dtoInterface := ctx.Locals("dto")
	if dtoInterface == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input: DTO is nil"})
	}

	dto, ok := dtoInterface.(*CreateUserDto)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error: invalid DTO type"})
	}

	user, err := c.service.CreateUser(ctx.Context(), *dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
// @Param user body UpdateUserDto true "Updated user information"
// @Success 200 {object} User
// @Failure 400 {object} FiberMap
// @Failure 500 {object} FiberMap
// @Router /api/v1/users/{id} [patch]
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	dtoInterface := ctx.Locals("dto")
	if dtoInterface == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input: DTO is nil"})
	}

	dto, ok := dtoInterface.(*UpdateUserDto)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error: invalid DTO type"})
	}

	user, err := c.service.UpdateUser(ctx.Context(), int32(id), *dto)
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
// @Router /api/v1/users/{id} [delete]
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

// SignIn handles user authentication and returns a JWT token
// @Summary User sign in
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body SignInInput true "User credentials"
// @Success 200 {object} SignInOutput
// @Failure 400 {object} FiberMap
// @Failure 401 {object} FiberMap
// @Router /api/v1/signin [post]
func (c *UserController) SignIn(ctx *fiber.Ctx) error {
	dtoInterface := ctx.Locals("dto")
	if dtoInterface == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input: DTO is nil"})
	}

	dto, ok := dtoInterface.(*SignInDto)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error: invalid DTO type"})
	}

	token, err := c.service.Authenticate(ctx.Context(), dto.Username, dto.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(SignInOutput{Token: token})
}

// SignOut handles user sign out (in this case, it's a client-side operation)
// @Summary User sign out
// @Description Sign out a user (client-side operation)
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} FiberMap
// @Router /api/v1/signout [post]
func (c *UserController) SignOut(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"message": "Successfully signed out"})
}
