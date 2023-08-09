package http

import (
	"fiberStore/dtos"
	"fiberStore/user/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
}

type userHandler struct {
	UserUsecase usecase.UserUsecase
	validator   *validator.Validate
}

func NewUserHandler(router *fiber.Group, UserUsecase usecase.UserUsecase, validator *validator.Validate) UserHandler {
	handler := &userHandler{
		UserUsecase: UserUsecase,
		validator:   validator,
	}

	// Main Routes
	api := router.Group("/user")

	// Routes
	api.Post("/register", handler.Register)

	return handler
}

func (user *userHandler) Register(c *fiber.Ctx) error {
	var (
		req dtos.UserRegister
		err error
	)

	if err = c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	if err = user.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error validating request body",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Context()
	result, err := user.UserUsecase.InsertOne(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse(
				fiber.StatusInternalServerError,
				"error registering user",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success registering user",
			result,
		),
	)
}
