//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/malytinKonstantin/go-fiber/internal/app"
	"github.com/malytinKonstantin/go-fiber/internal/db"
	"github.com/malytinKonstantin/go-fiber/internal/user"
)

var PostgresSet = wire.NewSet(
	db.NewPostgresPool,
	db.NewSQLDB,
)

var AppSet = wire.NewSet(
	PostgresSet,
	app.NewApp,
	user.NewModule,
	user.NewUserController,
	user.NewUserService,
	user.NewUserRepository,
)

func InitializeApp(databaseURL string) (*app.App, error) {
	wire.Build(AppSet)
	return nil, nil
}
