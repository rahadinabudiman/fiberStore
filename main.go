package main

import (
	"fiberStore/author"
	_cartRepository "fiberStore/cart/repository"
	_CartDetailHandler "fiberStore/cartDetail/delivery/http"
	_CartDetailRepository "fiberStore/cartDetail/repository"
	_CartDetailUsecase "fiberStore/cartDetail/usecase"
	_cloudinarUsecase "fiberStore/cloudinary/usecase"
	"fiberStore/helpers"
	"fiberStore/middlewares"
	_productHandler "fiberStore/product/delivery/http"
	_productRepository "fiberStore/product/repository"
	_productUsecase "fiberStore/product/usecase"
	_TransactionRepository "fiberStore/transaction/repository"
	_TransactionDetailHandler "fiberStore/transactionDetail/delivery/http"
	_TransactionDetailRepository "fiberStore/transactionDetail/repository"
	_TransactionDetailUsecase "fiberStore/transactionDetail/usecase"
	_userHandler "fiberStore/user/delivery/http"
	_userRepository "fiberStore/user/repository"
	_userUsecase "fiberStore/user/usecase"
	"fmt"
	"log"
	"os"
	"time"

	_ "fiberStore/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title           FiberStore Documentation API
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   r4ha
// @contact.url    https://github.com/rahadinabudiman

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// 54.179.176.114:1309/api/v1/swagger/index.html
// 54.179.176.114:1309

// @host      54.179.176.114:1309
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Setup Fiber
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
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

	// Seeding Data
	err = author.AccountSeeder(database)
	if err != nil {
		log.Fatal("Error seeding data")
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
	UserAmountUsecase := _userUsecase.NewUserAmountUsecase(UserAmountRepository, UserRepository, timeoutContext)
	UserUsecase := _userUsecase.NewUserUsecase(UserRepository, UserAmountRepository, timeoutContext)
	_userHandler.NewUserAmountHandler(admin.(*fiber.Group), UserAmountUsecase, UserUsecase, myValidator.GetValidator())
	_userHandler.NewUserHandler(api.(*fiber.Group), customer.(*fiber.Group), admin.(*fiber.Group), UserUsecase, myValidator.GetValidator())

	cloudinaryUsecase := _cloudinarUsecase.NewMediaUpload()
	productRepository := _productRepository.NewProductRepository(database)
	productUsecase := _productUsecase.NewProductUsecase(productRepository, UserRepository, timeoutContext)
	_productHandler.NewProductHandler(api.(*fiber.Group), admin.(*fiber.Group), productUsecase, cloudinaryUsecase, myValidator.GetValidator())

	CartRepository := _cartRepository.NewCartRepository(database)
	// CartUsecase := _cartUsecase.NewCartUsecase(CartRepository, timeoutContext)

	CartDetailRepository := _CartDetailRepository.NewCartDetailRepository(database)
	CartDetailUsecase := _CartDetailUsecase.NewCartDetailUsecase(CartDetailRepository, CartRepository, productRepository, UserRepository, timeoutContext)
	_CartDetailHandler.NewCartDetailHandler(api.(*fiber.Group), CartDetailUsecase, myValidator.GetValidator())

	TransactionRepository := _TransactionRepository.NewTransactionRepository(database)
	TransactionDetailRepository := _TransactionDetailRepository.NewTransactionDetailRepository(database)
	TransactionDetailUsecase := _TransactionDetailUsecase.NewTransactionDetailUsecase(TransactionDetailRepository, TransactionRepository, CartRepository, CartDetailRepository, productRepository, UserRepository, UserAmountRepository, timeoutContext)
	_TransactionDetailHandler.NewTransactionDetailHandler(api.(*fiber.Group), TransactionDetailUsecase, myValidator.GetValidator())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api.Get("/swagger/*", swagger.HandlerDefault)

	appPort := fmt.Sprintf(":%s", os.Getenv("SERVER_ADDRESS"))
	log.Fatal(app.Listen(appPort))
}
