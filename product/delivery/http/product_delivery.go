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
	FindOneProduct(c *fiber.Ctx) error
	FindAllByCategory(c *fiber.Ctx) error
	FindQueryAll(c *fiber.Ctx) error
	FindAll(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
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

	// Public Routes
	api.Get("/product/:id", handler.FindOneProduct)
	api.Get("/product", handler.FindQueryAll)
	api.Get("/product", handler.FindAll)
	api.Get("/product/", handler.FindAllByCategory)

	// Protected Admin Routes
	admin.Post("/product", handler.InsertProduct)
	admin.Put("/product/:id", handler.UpdateOne)
	admin.Delete("/product/:id", handler.DeleteOne)

	return handler
}

func (pd *productHandler) InsertProduct(c *fiber.Ctx) error {
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

	if err := pd.validator.Struct(req); err != nil {
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

	uploadUrl, err := pd.CloudinaryUsecase.FileUpload(models.File{File: src})
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

	res, err := pd.ProductUsecase.InsertOne(ctx, req, uploadUrl)
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

func (pd *productHandler) FindAllByCategory(c *fiber.Ctx) error {
	ctx := c.Context()

	pageParam := c.QueryInt("page")
	if pageParam == 0 {
		pageParam = 1
	}

	limitParam := c.QueryInt("limit")
	if limitParam == 0 {
		limitParam = 10
	}

	categoryParam := c.Query("category")

	result, count, err := pd.ProductUsecase.FindAllByCategory(ctx, pageParam, limitParam, categoryParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting all product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewPaginationResponse(
			fiber.StatusOK,
			"success getting all product",
			result,
			pageParam,
			limitParam,
			count,
		),
	)
}

func (pd *productHandler) FindOneProduct(c *fiber.Ctx) error {
	ProductID, err := strconv.Atoi(c.Params("id"))
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

	res, err := pd.ProductUsecase.FindOne(ctx, uint(ProductID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error finding product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success finding product",
			res,
		),
	)
}

func (pd *productHandler) FindQueryAll(c *fiber.Ctx) error {
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
	result, count, err := pd.ProductUsecase.FindQueryAll(ctx, pageParam, limitParam, searchParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting all product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewPaginationResponse(
			fiber.StatusOK,
			"success getting all product",
			result,
			pageParam,
			limitParam,
			count,
		),
	)
}

func (pd *productHandler) FindAll(c *fiber.Ctx) error {
	ctx := c.Context()

	pageParam := c.QueryInt("page")
	if pageParam == 0 {
		pageParam = 1
	}

	limitParam := c.QueryInt("limit")
	if limitParam == 0 {
		limitParam = 10
	}

	result, count, err := pd.ProductUsecase.FindAll(ctx, pageParam, limitParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error getting all product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewPaginationResponse(
			fiber.StatusOK,
			"success getting all product",
			result,
			pageParam,
			limitParam,
			count,
		),
	)
}

func (pd *productHandler) UpdateOne(c *fiber.Ctx) error {
	var req *dtos.UpdateProductRequest

	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(
			dtos.NewResponseMessage(
				fiber.StatusUnauthorized,
				"unauthorized",
			),
		)
	}

	AdministratorID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	ProductID, err := strconv.Atoi(c.Params("id"))
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

	req = &dtos.UpdateProductRequest{
		Name:     name,
		Detail:   detail,
		Price:    int64(price),
		Stock:    int64(stock),
		Category: category,
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error parsing request body",
				dtos.GetErrorData(err),
			),
		)
	}

	if err = pd.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error validating request body",
				dtos.GetErrorData(err),
			),
		)
	}

	product, err := pd.ProductUsecase.FindOne(ctx, uint(ProductID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"product not found",
				dtos.GetErrorData(err),
			),
		)
	}

	file, _ := c.FormFile("image")
	if file == nil {
		req.Image = product.Image
	} else {
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

		uploadUrl, err := pd.CloudinaryUsecase.FileUpload(models.File{File: src})
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				dtos.NewErrorResponse(
					fiber.StatusBadRequest,
					"Failed to upload file",
					dtos.GetErrorData(err),
				),
			)
		}

		req.Image = uploadUrl
	}

	productResponse, err := pd.ProductUsecase.UpdateOne(ctx, req, uint(ProductID), AdministratorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error updating product",
				dtos.GetErrorData(err),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		dtos.NewResponse(
			fiber.StatusOK,
			"success updating product",
			productResponse,
		),
	)
}

func (pd *productHandler) DeleteOne(c *fiber.Ctx) error {
	token := middlewares.GetTokenFromHeader(c)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewResponseMessage(
				fiber.StatusBadRequest,
				"unauthorized",
			),
		)
	}

	AdministratorID, err := middlewares.GetUserIdFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error get administratorID",
				dtos.GetErrorData(err),
			),
		)
	}

	ProductID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			dtos.NewErrorResponse(
				fiber.StatusBadRequest,
				"error get productID",
				dtos.GetErrorData(err),
			),
		)
	}

	ctx := c.Context()

	err = pd.ProductUsecase.DeleteOne(ctx, uint(ProductID), AdministratorID)
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
