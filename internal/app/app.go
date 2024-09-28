package app

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/malytinKonstantin/go-fiber/internal/user"
)

type App struct {
	UserModule *user.Module
	DB         *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		UserModule: user.SetupModule(db),
		DB:         db,
	}
}

func (a *App) SetupRoutes(app *fiber.App) {
	a.UserModule.SetupRoutes(app)
}
