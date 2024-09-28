package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app, err := InitializeApp(db)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	fiberApp := fiber.New()
	app.SetupRoutes(fiberApp)

	log.Fatal(fiberApp.Listen(":3000"))
}
