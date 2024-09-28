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

func (a *App) GetService(key shared.ModuleKey) (any, bool) {
	for _, module := range a.Modules {
		if service := module.GetService(key); service != nil {
			return service, true
		}
	}
	return nil, false
}

func (a *App) GetRepository(key shared.ModuleKey) (any, bool) {
	for _, module := range a.Modules {
		if repository := module.GetRepository(key); repository != nil {
			return repository, true
		}
	}
	return nil, false
}
