package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CartDetailHandler interface {
	InsertOne(c *fiber.Ctx) error
	GetCartDetail(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type cartDetailHandler struct {
	CartDetailUsecase models.CartDetailUsecase
	validator         *validator.Validate
}

func NewCartDetailHandler(protected *fiber.Group, CartDetailUsecase models.CartDetailUsecase, validator *validator.Validate) CartDetailHandler {
	handler := &cartDetailHandler{
		CartDetailUsecase: CartDetailUsecase,
		validator:         validator,
	}

	// Protected User Routes
	protected.Post("/cart", handler.InsertOne)
	protected.Get("/cart", handler.GetCartDetail)
	protected.Delete("/cart", handler.DeleteProduct)

	return handler
}

func (ch *cartDetailHandler) InsertOne(c *fiber.Ctx) error {
	var req dtos.InsertCartDetailRequest

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

	result, err := ch.CartDetailUsecase.InsertOne(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error inserting CartDetail",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		dtos.NewResponse(
			fiber.StatusCreated,
			"success inserting CartDetail",
			result,
		),
	)
}

func (ch *cartDetailHandler) GetCartDetail(c *fiber.Ctx) error {
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

	result, err := ch.CartDetailUsecase.FindAll(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting CartDetails",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success getting CartDetails",
			result,
		),
	)
}

func (ch *cartDetailHandler) DeleteProduct(c *fiber.Ctx) error {
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

	productID := c.QueryInt("product_id")

	ctx := c.Context()

	err = ch.CartDetailUsecase.DeleteProduct(ctx, userID, uint(productID))
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
