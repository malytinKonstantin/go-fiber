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

func NewApp(userModule *user.Module, db *sql.DB) *App {
	return &App{
		UserModule: userModule,
		DB:         db,
	}
}

func (a *App) SetupRoutes(router fiber.Router) {
	a.UserModule.SetupRoutes(router)
}
