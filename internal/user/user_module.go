package user

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	Controller *UserController
}

func SetupModule(db *sql.DB) *Module {
	userRepo := NewUserRepository(db)
	userService := NewUserService(userRepo)
	userController := NewUserController(userService)

	return &Module{
		Controller: userController,
	}
}

func (m *Module) SetupRoutes(app *fiber.App) {
	SetupRoutes(app, m.Controller)
}
