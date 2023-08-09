package main

import (
	"fiberStore/author"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Setup Fiber
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	// Initialize default config for logger middleware
	app.Use(logger.New())

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Setup Configuration
	database, err := author.InitMySQL()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	// Migrates Database
	err = author.MigrateDB(database)
	if err != nil {
		log.Fatal("Error migrating database")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// api := app.Group("/api/v1")

	appPort := fmt.Sprintf(":%s", os.Getenv("SERVER_ADDRESS"))
	log.Fatal(app.Listen(appPort))

}
