package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/malytinKonstantin/go-fiber/internal/app"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация приложения
	application := app.NewApp(db)

	// Создание Fiber приложения
	fiberApp := fiber.New()

	// Настройка маршрутов
	application.SetupRoutes(fiberApp)

	// Запуск сервера
	log.Fatal(fiberApp.Listen(":3000"))
}
