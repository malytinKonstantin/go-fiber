package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
	_ "github.com/malytinKonstantin/go-fiber/docs"
	"github.com/spf13/viper"
)

func main() {
	dbURL := viper.GetString("DATABASE_URL")
	port := viper.GetString("PORT")
	apiPrefix := viper.GetString("API_PREFIX")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app, err := InitializeApp(db)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	fiberApp := fiber.New(fiber.Config{
		JSONDecoder: json.Unmarshal,
	})
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)
	api := fiberApp.Group(apiPrefix)
	app.SetupRoutes(api)

	log.Fatal(fiberApp.Listen(fmt.Sprintf(":%s", port)))
}

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
	}
}
