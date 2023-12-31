package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransactionDetailHandler interface {
	InsertOne(c *fiber.Ctx) error
	FindLatestTransaction(c *fiber.Ctx) error
}

type transactionDetailHandler struct {
	TransactionDetailUsecase models.TransactionDetailUsecase
	validator                *validator.Validate
}

func NewTransactionDetailHandler(protected *fiber.Group, TransactionDetailUsecase models.TransactionDetailUsecase, validator *validator.Validate) TransactionDetailHandler {
	handler := &transactionDetailHandler{
		TransactionDetailUsecase: TransactionDetailUsecase,
		validator:                validator,
	}

	// Protected User Routes
	protected.Post("/transaction", handler.InsertOne)
	protected.Get("/transaction", handler.FindLatestTransaction)

	return handler
}

func (th *transactionDetailHandler) InsertOne(c *fiber.Ctx) error {
	var req models.TransactionDetail

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

	req.UserID = userID

	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"bad request",
				dtos.GetErrorData(err),
			),
		)
	}

	err = th.validator.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"bad request",
				dtos.GetErrorData(err),
			),
		)
	}

	res, err := th.TransactionDetailUsecase.InsertOne(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse(
				fiber.StatusInternalServerError,
				"internal server error",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		dtos.NewResponse(
			fiber.StatusCreated,
			"success",
			res,
		),
	)
}

func (th *transactionDetailHandler) FindLatestTransaction(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"bad request",
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

	res, err := th.TransactionDetailUsecase.FindAll(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse(
				fiber.StatusInternalServerError,
				"internal server error",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success",
			res,
		),
	)
}
