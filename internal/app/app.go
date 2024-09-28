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
func (a *App) GetService(serviceName string) shared.Service {
	for _, module := range a.Modules {
		if service := module.GetService(serviceName); service != nil {
			return service
		}
	}
	return nil
}

// Метод для получения репозитория из определенного модуля
func (a *App) GetRepository(repositoryName string) shared.Repository {
	for _, module := range a.Modules {
		if repository := module.GetRepository(repositoryName); repository != nil {
			return repository
		}
	}
	return nil
}
