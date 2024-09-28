package user

import (
	"github.com/gofiber/fiber/v2"
)

type Module struct {
	Controller *UserController
}

func NewModule(controller *UserController) *Module {
	return &Module{
		Controller: controller,
	}
}

func (m *Module) SetupRoutes(router fiber.Router) {
	m.Controller.SetupRoutes(router)
}
