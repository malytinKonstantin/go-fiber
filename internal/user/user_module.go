package user

import (
	"database/sql"

	"github.com/malytinKonstantin/go-fiber/internal/shared"
)

const (
	UserServiceKey    shared.ModuleKey = "user_service"
	UserRepositoryKey shared.ModuleKey = "user_repository"
)

func SetupModule(db *sql.DB) *shared.Module {
	userRepo := NewUserRepository(db)
	userService := NewUserService(userRepo)
	userController := NewUserController(userService)

	module := shared.NewModule(userController)

	module.AddRepository(UserRepositoryKey, userRepo)
	module.AddService(UserServiceKey, userService)

	return module
}
