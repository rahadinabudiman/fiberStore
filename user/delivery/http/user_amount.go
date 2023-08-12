package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserAmountHandler interface {
	TopUpSaldo(c *fiber.Ctx) error
}

type userAmountHandler struct {
	UserAmountUsecase models.UserAmountUsecase
	UserUsecase       models.UserUsecase
	validator         *validator.Validate
}

func NewUserAmountHandler(admin *fiber.Group, UserAmountUsecase models.UserAmountUsecase, UserUsecase models.UserUsecase, validator *validator.Validate) UserAmountHandler {
	handler := &userAmountHandler{
		UserAmountUsecase: UserAmountUsecase,
		UserUsecase:       UserUsecase,
		validator:         validator,
	}

	// Route
	admin.Post("/topup", handler.TopUpSaldo)

	return handler
}

func (uah *userAmountHandler) TopUpSaldo(c *fiber.Ctx) error {
	var TopUp *dtos.TopUpSaldoRequest

	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"error getting token from header",
			),
		)
	}

	role, err := middlewares.GetRoleFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting role from token",
				dtos.GetErrorData(err),
			),
		)
	}

	if role != "Admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewResponseMessage(
				fiber.StatusUnauthorized,
				"Unauhorized",
			),
		)
	}

	if err = c.BodyParser(&TopUp); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			dtos.NewErrorResponse(
				fiber.StatusUnprocessableEntity,
				"error parsing body",
				dtos.GetErrorData(err),
			),
		)
	}

	if err = uah.validator.Struct(TopUp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error validating request body",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Context()

	result, err := uah.UserAmountUsecase.TopUpSaldo(ctx, TopUp)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"failed to top up balance",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"Success Top Up Balance",
			result,
		),
	)
}
