package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler interface {
	// Authentikasi User
	LoginUser(c *fiber.Ctx) error
	// CRUD User
	Register(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
	GetAllProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UpdateProfilePassword(c *fiber.Ctx) error
	DeleteAccount(c *fiber.Ctx) error
}

type userHandler struct {
	UserUsecase models.UserUsecase
	validator   *validator.Validate
}

func NewUserHandler(api *fiber.Group, user *fiber.Group, admin *fiber.Group, UserUsecase models.UserUsecase, validator *validator.Validate) UserHandler {
	handler := &userHandler{
		UserUsecase: UserUsecase,
		validator:   validator,
	}
	// Authentikasi
	api.Post("/login", handler.LoginUser)
	api.Post("/register", handler.Register)

	// Protected User Routes
	user.Get("", handler.GetProfile)
	user.Put("", handler.UpdateProfile)
	user.Put("/password", handler.UpdateProfilePassword)
	user.Delete("", handler.DeleteAccount)

	// Protected Admin Routes
	admin.Get("/user", handler.GetAllProfile)
	admin.Delete("/user", handler.DeleteAccount)

	return handler
}

func (user *userHandler) LoginUser(c *fiber.Ctx) error {
	var (
		req dtos.UserLoginRequest
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

	result, err := user.UserUsecase.LoginUser(ctx, c, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error logging in user",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success logging in user",
			result,
		),
	)
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
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error registering user",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success registering user",
			result,
		),
	)
}

func (user *userHandler) GetProfile(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"error getting token from header",
			),
		)
	}

	UserID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"error getting user id from token",
			),
		)
	}

	ctx := c.Context()

	result, err := user.UserUsecase.FindOneById(ctx, int(UserID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			dtos.NewErrorResponse(
				fiber.StatusInternalServerError,
				"error getting user profile",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success getting user profile",
			result,
		),
	)
}

func (user *userHandler) GetAllProfile(c *fiber.Ctx) error {
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
				"error unauthorized",
			),
		)
	}

	ctx := c.Context()

	pageParam := c.QueryInt("page")
	if pageParam == 0 {
		pageParam = 1
	}

	limitParam := c.QueryInt("limit")
	if limitParam == 0 {
		limitParam = 10
	}

	searchParam := c.Query("search")
	sortByParam := c.Query("sortBy")
	if sortByParam == "" {
		sortByParam = "asc"
	}

	result, count, err := user.UserUsecase.FindAll(ctx, pageParam, limitParam, searchParam, sortByParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting all user profile",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewPaginationResponse(
			fiber.StatusOK,
			"success getting all user profile",
			result,
			pageParam,
			limitParam,
			count,
		),
	)
}

func (user *userHandler) UpdateProfile(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"error getting token from header",
			),
		)
	}

	userID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting user id from token",
				dtos.GetErrorData(err),
			),
		)
	}

	var req dtos.UserUpdateProfile

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

	result, err := user.UserUsecase.UpdateOne(ctx, int(userID), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error updating user profile",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success updating user profile",
			result,
		),
	)
}

func (user *userHandler) UpdateProfilePassword(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"error getting token from header",
			),
		)
	}

	userID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting user id from token",
				dtos.GetErrorData(err),
			),
		)
	}

	var req dtos.UserUpdatePassword

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

	result, err := user.UserUsecase.UpdatePassword(ctx, uint(userID), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error updating user profile password",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success updating password",
			result,
		),
	)
}

func (user *userHandler) DeleteAccount(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"error getting token from header",
			),
		)
	}

	userID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting user id from token",
				dtos.GetErrorData(err),
			),
		)
	}

	var req dtos.DeleteUserRequest

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

	err = user.UserUsecase.DeleteOne(ctx, uint(userID), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error deleting user profile",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponseMessage(
			fiber.StatusOK,
			"success deleting user profile",
		),
	)
}
