package user

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/middleware"
)

const (
	ErrInvalidID          = "Invalid ID"
	ErrUserNotFound       = "User not found"
	ErrInvalidDTO         = "Invalid input: DTO is nil"
	ErrInvalidDTOType     = "Internal server error: invalid DTO type"
	ErrFailedToUpdateUser = "Failed to update user"
	ErrFailedToDeleteUser = "Failed to delete user"
	ErrInvalidQueryParams = "Invalid query parameters"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func sendErrorResponse(ctx *fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(fiber.Map{"error": message})
}

func getDTO[T any](ctx *fiber.Ctx) (*T, error) {
	dtoInterface := ctx.Locals("dto")
	if dtoInterface == nil {
		return nil, errors.New(ErrInvalidDTO)
	}
	dto, ok := dtoInterface.(*T)
	if !ok {
		return nil, errors.New(ErrInvalidDTOType)
	}
	return dto, nil
}

// SetupRoutes sets up the user-related routes
// @Summary Set up user routes
// @Description Set up routes for user-related operations
// @Tags users
func (c *UserController) SetupRoutes(router fiber.Router) {
	// public routes
	router.Post("/signin", middleware.SkipAuth(c.SignIn))
	middleware.RegisterDTO("/signin", "POST", SignInDto{})
	router.Post("/signup", middleware.SkipAuth(c.CreateUser))

	// protected routes
	router.Get("/users", c.ListUsers)
	router.Get("/users/:id", c.GetUser)
	router.Get("/users/username/:username", c.GetUserByUsername)
	middleware.RegisterDTO("/users", "POST", CreateUserDto{})
	router.Post("/users", c.CreateUser)
	middleware.RegisterDTO("/users/:id", "PATCH", UpdateUserDto{})
	router.Patch("/users/:id", c.UpdateUser)
	router.Delete("/users/:id", c.DeleteUser)
}

// GetUser retrieves a user by ID
// @Summary Get a user
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400,404 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, ErrInvalidID)
	}

	user, err := c.service.GetUser(ctx.Context(), int32(id))
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusNotFound, ErrUserNotFound)
	}

	return ctx.JSON(user)
}

// GetUserByUsername retrieves a user by username
// @Summary Get a user by username
// @Tags users
// @Param username path string true "Username"
// @Success 200 {object} User
// @Failure 400,404 {object} ErrorResponse
// @Router /api/v1/users/username/{username} [get]
func (c *UserController) GetUserByUsername(ctx *fiber.Ctx) error {
	username := ctx.Params("username")
	user, err := c.service.GetUserByUsername(ctx.Context(), username)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusNotFound, ErrUserNotFound)
	}

	return ctx.JSON(user)
}

// ListUsers retrieves a list of users based on query parameters
// @Summary List users
// @Tags users
// @Param username query string false "Username"
// @Param email query string false "Email"
// @Param full_name query string false "Full Name"
// @Param bio query string false "Bio"
// @Param created_from query string false "Created From (YYYY-MM-DD)"
// @Param created_to query string false "Created To (YYYY-MM-DD)"
// @Param sort_by query string false "Sort By (e.g., username_asc, created_at_desc)"
// @Param limit query int false "Limit" default(100)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} User
// @Failure 400,500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	query := new(ListUsersQuery)
	if err := ctx.QueryParser(query); err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, ErrInvalidQueryParams)
	}

	params := SearchUsersParams{
		Username:    query.Username,
		Email:       query.Email,
		FullName:    query.FullName,
		Bio:         query.Bio,
		CreatedFrom: query.CreatedFrom,
		CreatedTo:   query.CreatedTo,
		SortBy:      query.SortBy,
		Limit:       query.Limit,
		Offset:      query.Offset,
	}

	if params.Limit <= 0 {
		params.Limit = 100
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	users, err := c.service.SearchUsers(ctx.Context(), params)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(users)
}

// CreateUser creates a new user
// @Summary Create a user
// @Tags users
// @Param user body CreateUserDto true "User information"
// @Success 201 {object} User
// @Failure 400,500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	dto, err := getDTO[CreateUserDto](ctx)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	user, err := c.service.CreateUser(ctx.Context(), *dto)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Tags users
// @Param id path int true "User ID"
// @Param user body UpdateUserDto true "Updated user information"
// @Success 200 {object} User
// @Failure 400,500 {object} ErrorResponse
// @Router /api/v1/users/{id} [patch]
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, ErrInvalidID)
	}

	dto, err := getDTO[UpdateUserDto](ctx)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	user, err := c.service.UpdateUser(ctx.Context(), int32(id), *dto)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusInternalServerError, ErrFailedToUpdateUser)
	}

	return ctx.JSON(user)
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Tags users
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400,500 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, ErrInvalidID)
	}

	if err := c.service.DeleteUser(ctx.Context(), int32(id)); err != nil {
		return sendErrorResponse(ctx, fiber.StatusInternalServerError, ErrFailedToDeleteUser)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// SignIn handles user authentication and returns a JWT token
// @Summary User sign in
// @Tags auth
// @Param credentials body SignInDto true "User credentials"
// @Success 200 {object} SignInOutput
// @Failure 400,401 {object} ErrorResponse
// @Router /api/v1/signin [post]
func (c *UserController) SignIn(ctx *fiber.Ctx) error {
	dto, err := getDTO[SignInDto](ctx)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	token, err := c.service.Authenticate(ctx.Context(), dto.Username, dto.Password)
	if err != nil {
		return sendErrorResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}

	return ctx.JSON(SignInOutput{Token: token})
}

// SignOut handles user sign out (client-side operation)
// @Summary User sign out
// @Tags auth
// @Success 200 {object} SuccessResponse
// @Router /api/v1/signout [post]
func (c *UserController) SignOut(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"message": "Successfully signed out"})
}
