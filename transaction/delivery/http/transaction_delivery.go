package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler interface {
	InsertOne(c *fiber.Ctx) error
}

type transactionHandler struct {
	TransactionUsecase models.TransactionUsecase
	validator          *validator.Validate
}

func NewTransactionHandler(protected *fiber.Group, TransactionUsecase models.TransactionUsecase, validator *validator.Validate) TransactionHandler {
	handler := &transactionHandler{
		TransactionUsecase: TransactionUsecase,
		validator:          validator,
	}

	// Protected User Routes
	protected.Post("/transaction", handler.InsertOne)

	return handler
}

func (th *transactionHandler) InsertOne(c *fiber.Ctx) error {
	var req models.Transaction

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

	res, err := th.TransactionUsecase.InsertOne(c.Context(), &req)
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
