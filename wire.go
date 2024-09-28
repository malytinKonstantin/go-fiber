//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/malytinKonstantin/go-fiber/internal/app"
	"github.com/malytinKonstantin/go-fiber/internal/user"
)

func InitializeApp(db *sql.DB) (*app.App, error) {
	wire.Build(
		app.NewApp,
		user.NewModule,
		user.NewUserController,
		user.NewUserService,
		user.NewUserRepository,
	)
	return &app.App{}, nil
}
