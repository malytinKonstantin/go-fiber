package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/malytinKonstantin/sqlc-test/internal/controllers"
	"github.com/malytinKonstantin/sqlc-test/internal/repositories"
	"github.com/malytinKonstantin/sqlc-test/internal/services"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация репозитория, сервиса и контроллера
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Создание Fiber приложения
	app := fiber.New()

	// Определение маршрутов
	app.Get("/users", userController.ListUsers)
	app.Get("/users/:id", userController.GetUser)
	app.Post("/users", userController.CreateUser)
	app.Delete("/users/:id", userController.DeleteUser)

	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}
