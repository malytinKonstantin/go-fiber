package user

import (
	"database/sql"

	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

func SetupModule(db *sql.DB) *shared.Module {
	userRepo := NewUserRepository(db)
	userService := NewUserService(userRepo)
	userController := NewUserController(userService)

	module := shared.NewModule(userController)

	module.AddRepository("user", userRepo)
	module.AddService("user", userService)

	// Добавьте дополнительные сервисы и репозитории, если они есть
	// Например:
	// authService := NewAuthService(userRepo)
	// module.AddService("auth", authService)

	return module
}
