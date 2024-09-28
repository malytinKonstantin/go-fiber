package shared

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	SetupRoutes(app *fiber.App)
}

type Service interface {
	// Общий интерфейс для сервисов, если необходимо
}

type Repository interface {
	// Общий интерфейс для репозиториев, если необходимо
}

type Module struct {
	Controller   Controller
	Services     map[string]Service
	Repositories map[string]Repository
}

func NewModule(controller Controller) *Module {
	return &Module{
		Controller:   controller,
		Services:     make(map[string]Service),
		Repositories: make(map[string]Repository),
	}
}

func (m *Module) SetupRoutes(app *fiber.App) {
	m.Controller.SetupRoutes(app)
}

func (m *Module) AddService(name string, service Service) {
	m.Services[name] = service
}

func (m *Module) GetService(name string) Service {
	return m.Services[name]
}

func (m *Module) AddRepository(name string, repository Repository) {
	m.Repositories[name] = repository
}

func (m *Module) GetRepository(name string) Repository {
	return m.Repositories[name]
}
