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

type ModuleKey string

const (
	ServiceKey    ModuleKey = "service"
	RepositoryKey ModuleKey = "repository"
)

type Module struct {
	Controller   Controller
	Services     map[ModuleKey]Service
	Repositories map[ModuleKey]Repository
}

func NewModule(controller Controller) *Module {
	return &Module{
		Controller:   controller,
		Services:     make(map[ModuleKey]Service),
		Repositories: make(map[ModuleKey]Repository),
	}
}

func (m *Module) SetupRoutes(app *fiber.App) {
	m.Controller.SetupRoutes(app)
}

func (m *Module) AddService(key ModuleKey, service Service) {
	m.Services[key] = service
}

func (m *Module) GetService(key ModuleKey) Service {
	return m.Services[key]
}

func (m *Module) AddRepository(key ModuleKey, repository Repository) {
	m.Repositories[key] = repository
}

func (m *Module) GetRepository(key ModuleKey) Repository {
	return m.Repositories[key]
}
