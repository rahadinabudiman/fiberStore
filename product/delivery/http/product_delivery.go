package http

import (
	"fiberStore/dtos"
	"fiberStore/middlewares"
	"fiberStore/models"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler interface {
	InsertProduct(c *fiber.Ctx) error
}

type productHandler struct {
	ProductUsecase    models.ProductUsecase
	CloudinaryUsecase models.CloudinaryUsecase
	validator         *validator.Validate
}

func NewProductHandler(api *fiber.Group, admin *fiber.Group, ProductUsecase models.ProductUsecase, CloudinaryUsecase models.CloudinaryUsecase, validator *validator.Validate) ProductHandler {
	handler := &productHandler{
		ProductUsecase:    ProductUsecase,
		CloudinaryUsecase: CloudinaryUsecase,
		validator:         validator,
	}

	// Protected Admin Routes
	admin.Post("/product", handler.InsertProduct)

	return handler
}

func (product *productHandler) InsertProduct(c *fiber.Ctx) error {
	var req *dtos.InserProductRequest

	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewResponseMessage(
				fiber.StatusUnauthorized,
				"unauthorized",
			),
		)
	}

	UserID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Context()

	name := c.FormValue("name")
	detail := c.FormValue("detail")
	price, _ := strconv.Atoi(c.FormValue("price"))
	stock, _ := strconv.Atoi(c.FormValue("stock"))
	category := c.FormValue("category")

	req = &dtos.InserProductRequest{
		AdministratorID: UserID,
		Name:            name,
		Detail:          detail,
		Price:           int64(price),
		Stock:           int64(stock),
		Category:        category,
	}

	// Validate
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	if err := product.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error validating request body",
				dtos.GetErrorData(err),
			),
		)
	}

	// Get Image
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"Failed to open file",
				dtos.GetErrorData(err),
			),
		)
	}
	defer src.Close()

	// Validate extension images
	re := regexp.MustCompile(`.png|.jpeg|.jpg`)
	if !re.MatchString(file.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"The provided file format is not allowed. Please upload a JPEG or PNG image",
				dtos.GetErrorData(err),
			),
		)
	}

	uploadUrl, err := product.CloudinaryUsecase.FileUpload(models.File{File: src})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"Failed to upload file",
				dtos.GetErrorData(err),
			),
		)
	}

	req.Image = file

	res, err := product.ProductUsecase.InsertOne(ctx, req, uploadUrl)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error inserting product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(
		dtos.NewResponse(
			fiber.StatusCreated,
			"success inserting product",
			res,
		),
	)
}
