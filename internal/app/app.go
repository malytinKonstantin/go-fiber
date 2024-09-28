package app

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/shared"
	"github.com/malytinKonstantin/go-fiber/internal/user"
)

type App struct {
	Modules []*shared.Module
	DB      *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		Modules: []*shared.Module{
			user.SetupModule(db),
			// Добавьте другие модули здесь
		},
		DB: db,
	}
}

func (a *App) SetupRoutes(app *fiber.App) {
	for _, module := range a.Modules {
		module.SetupRoutes(app)
	}
}

// Метод для получения сервиса из определенного модуля
func (a *App) GetService(moduleName string, serviceName string) shared.Service {
	for _, module := range a.Modules {
		if controller, ok := module.Controller.(interface{ ModuleName() string }); ok {
			if controller.ModuleName() == moduleName {
				return module.GetService(serviceName)
			}
		}
	}
	return nil
}

// Метод для получения репозитория из определенного модуля
func (a *App) GetRepository(moduleName string, repositoryName string) shared.Repository {
	for _, module := range a.Modules {
		if controller, ok := module.Controller.(interface{ ModuleName() string }); ok {
			if controller.ModuleName() == moduleName {
				return module.GetRepository(repositoryName)
			}
		}
	}
	return nil
}
