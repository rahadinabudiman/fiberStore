package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CartHandler interface {
	InsertOne(c *fiber.Ctx) error
	GetCart(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type cartHandler struct {
	CartUsecase models.CartUsecase
	validator   *validator.Validate
}

func NewCartHandler(protected *fiber.Group, CartUsecase models.CartUsecase, validator *validator.Validate) CartHandler {
	handler := &cartHandler{
		CartUsecase: CartUsecase,
		validator:   validator,
	}

	// Protected User Routes
	protected.Post("/cart", handler.InsertOne)
	protected.Get("/cart", handler.GetCart)
	protected.Delete("/cart/:product_id", handler.DeleteProduct)

	return handler
}

func (ch *cartHandler) InsertOne(c *fiber.Ctx) error {
	var req dtos.InsertCartRequest

	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewResponseMessage(
				fiber.StatusUnauthorized,
				"unauthorized",
			),
		)
	}

	userID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewErrorResponse(
				fiber.StatusUnauthorized,
				"unauthorized",
				dtos.GetErrorData(err),
			),
		)
	}

	if err = c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	if err = ch.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error validating request body",
				dtos.GetErrorData(err),
			),
		)
	}

	req.UserID = userID

	ctx := c.Context()

	result, err := ch.CartUsecase.InsertOne(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error inserting cart",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		dtos.NewResponse(
			fiber.StatusCreated,
			"success inserting cart",
			result,
		),
	)
}

func (ch *cartHandler) GetCart(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewResponseMessage(
				fiber.StatusUnauthorized,
				"unauthorized",
			),
		)
	}

	userID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewErrorResponse(
				fiber.StatusUnauthorized,
				"unauthorized",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Context()

	result, err := ch.CartUsecase.FindAll(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting carts",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success getting carts",
			result,
		),
	)
}

func (ch *cartHandler) DeleteProduct(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewResponseMessage(
				fiber.StatusUnauthorized,
				"unauthorized",
			),
		)
	}

	userID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewErrorResponse(
				fiber.StatusUnauthorized,
				"unauthorized",
				dtos.GetErrorData(err),
			),
		)
	}

	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error converting product id",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Context()

	err = ch.CartUsecase.DeleteProduct(ctx, userID, uint(productID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error deleting product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponseMessage(
			fiber.StatusOK,
			"success deleting product",
		),
	)
}
