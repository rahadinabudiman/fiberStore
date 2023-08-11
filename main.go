package main

import (
	"fiberStore/author"
	"fiberStore/helpers"
	"fiberStore/middlewares"
	_userHandler "fiberStore/user/delivery/http"
	_userRepository "fiberStore/user/repository"
	_userUsecase "fiberStore/user/usecase"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
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
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_KEY"),
	}))

	api := app.Group("/api/v1")
	customer := api.Group("/user")
	admin := api.Group("/admin")

	myValidator := helpers.NewXValidator()

	// Setup Security Routes
	customer.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("Customer"))
	admin.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("Admin"))

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

	// Setup Context Timeout
	CONTEXT_TIMEOUT, err := helpers.GetEnvInt("CONTEXT_TIMEOUT")
	if err != nil {
		log.Fatal(err)
	}
	timeoutContext := time.Duration(CONTEXT_TIMEOUT) * time.Second

	// Setup Routes
	UserAmountRepository := _userRepository.NewUserAmountRepository(database)
	UserRepository := _userRepository.NewUserRepository(database)

	// UserAmountUsecase := _userUsecase.NewUserAmountUsecase(UserAmountRepository, UserRepository, timeoutContext)

	UserUsecase := _userUsecase.NewUserUsecase(UserRepository, UserAmountRepository, timeoutContext)
	_userHandler.NewUserHandler(api.(*fiber.Group), customer.(*fiber.Group), admin.(*fiber.Group), UserUsecase, myValidator.GetValidator())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	appPort := fmt.Sprintf(":%s", os.Getenv("SERVER_ADDRESS"))
	log.Fatal(app.Listen(appPort))
}
