package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/helpers"
	"fiberStore/models"
	"time"
)

type productUsecase struct {
	ProductRepository models.ProductRepository
	UserRepository    models.UserRepository
	contextTimeout    time.Duration
}

func NewProductUsecase(ProductRepository models.ProductRepository, UserRepository models.UserRepository, contextTimeout time.Duration) models.ProductUsecase {
	return &productUsecase{
		ProductRepository: ProductRepository,
		UserRepository:    UserRepository,
		contextTimeout:    contextTimeout,
	}
}

// AddProduct godoc
// @Summary Add Product
// @Description Add Product
// @Tags Admin - Product
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param detail formData string true "Product Detail"
// @Param price formData int64 true "Product Price"
// @Param stock formData int64 true "Product Stock"
// @Param category formData string true "Product Category"
// @Param image formData file false "Photo Product"
// @Success 200 {object} dtos.InsertProductStatusOKResponse
// @Failure 400 {object} dtos.BadRequestResponse
// @Failure 401 {object} dtos.UnauthorizedResponse
// @Failure 403 {object} dtos.ForbiddenResponse
// @Failure 404 {object} dtos.NotFoundResponse
// @Failure 500 {object} dtos.InternalServerErrorResponse
// @Router /admin/product [post]
// @Security BearerAuth
func (pu *productUsecase) InsertOne(ctx context.Context, req *dtos.InserProductRequest, url string) (*dtos.InserProductResponse, error) {
	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Check Admin
	_, err := pu.UserRepository.FindOneAdmin(req.AdministratorID)
	if err != nil {
		return nil, errors.New("not allowed to create product")
	}

	slug := helpers.CreateSlug(req.Name)
	imageUrl := url

	// Check Slug Exist
	checkSlug, err := pu.ProductRepository.FindOneBySlug(slug)
	if err == nil {
		if checkSlug.Slug == slug {
			slug = slug + "-" + helpers.RandomString(5)
		}
	}

	CreateProduct := &models.Product{
		AdministratorID: req.AdministratorID,
		Slug:            slug,
		Name:            req.Name,
		Detail:          req.Detail,
		Price:           req.Price,
		Stock:           req.Stock,
		Category:        req.Category,
		Image:           imageUrl,
	}

	createdProduct, err := pu.ProductRepository.InsertOne(CreateProduct)
	if err != nil {
		return nil, errors.New("failed to create product")
	}

	res := &dtos.InserProductResponse{
		Name:     createdProduct.Name,
		Slug:     createdProduct.Slug,
		Detail:   createdProduct.Detail,
		Price:    createdProduct.Price,
		Stock:    createdProduct.Stock,
		Category: createdProduct.Category,
		Image:    createdProduct.Image,
	}

	return res, nil
}

// GetAllProduct godoc
// @Summary      Get all products
// @Description  Get all products
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllProductStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /product [get]
func (pu *productUsecase) FindAll(ctx context.Context, page, limit int) (*[]dtos.ProductResponse, int, error) {
	var productResponses []dtos.ProductResponse

	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	products, count, err := pu.ProductRepository.FindAll(page, limit)
	if err != nil {
		return nil, 0, errors.New("error getting products")
	}

	for _, product := range *products {
		productResponses = append(productResponses, dtos.ProductResponse{
			ID:       product.ID,
			Slug:     product.Slug,
			Name:     product.Name,
			Detail:   product.Detail,
			Price:    product.Price,
			Stock:    product.Stock,
			Category: product.Category,
			Image:    product.Image,
		})
	}

	return &productResponses, count, nil
}

// GetAllProductByName godoc
// @Summary      Get all products by name
// @Description  Get all products by name
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param name query string false "Product name"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllProductStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /product/findByName [get]
func (pu *productUsecase) FindQueryAll(ctx context.Context, page, limit int, search string) (*[]dtos.ProductResponse, int, error) {
	var productResponses []dtos.ProductResponse

	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	products, count, err := pu.ProductRepository.FindQueryAll(page, limit, search)
	if err != nil {
		return nil, 0, errors.New("error getting products")
	}

	for _, product := range *products {
		productResponse := dtos.ProductResponse{
			ID:       product.ID,
			Slug:     product.Slug,
			Name:     product.Name,
			Detail:   product.Detail,
			Price:    product.Price,
			Stock:    product.Stock,
			Category: product.Category,
			Image:    product.Image,
		}
		productResponses = append(productResponses, productResponse)
	}

	return &productResponses, count, nil
}

// GetAllProductByCategory godoc
// @Summary      Get all products by Category
// @Description  Get all products by Category
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param category query string false "Product Category"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllProductStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /product/findByCategory [get]
func (pu *productUsecase) FindAllByCategory(ctx context.Context, page, limit int, category string) (*[]dtos.ProductResponse, int, error) {
	var productResponses []dtos.ProductResponse

	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	products, count, err := pu.ProductRepository.FindByCategory(page, limit, category)
	if err != nil {
		return nil, 0, errors.New("error getting products")
	}

	for _, product := range *products {
		productResponse := dtos.ProductResponse{
			ID:       product.ID,
			Slug:     product.Slug,
			Name:     product.Name,
			Detail:   product.Detail,
			Price:    product.Price,
			Stock:    product.Stock,
			Category: product.Category,
			Image:    product.Image,
		}
		productResponses = append(productResponses, productResponse)
	}

	return &productResponses, count, nil
}

func (pu *productUsecase) FindOne(ctx context.Context, id uint) (*dtos.ProductResponse, error) {
	var productResponse dtos.ProductResponse

	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	product, err := pu.ProductRepository.FindOne(int(id))
	if err != nil {
		return nil, errors.New("error getting product")
	}

	productResponse = dtos.ProductResponse{
		ID:       product.ID,
		Slug:     product.Slug,
		Name:     product.Name,
		Detail:   product.Detail,
		Price:    product.Price,
		Stock:    product.Stock,
		Category: product.Category,
		Image:    product.Image,
	}

	return &productResponse, nil
}

// UpdateProduct godoc
// @Summary      Update Product
// @Description  Update Product
// @Tags         Admin - Product
// @Accept       json
// @Produce      json
// @Param id query int false "Product ID"
// @Param        request body dtos.UpdateProductRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ProductStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/product/{id} [put]
// @Security BearerAuth
func (pu *productUsecase) UpdateOne(ctx context.Context, req *dtos.UpdateProductRequest, id, AdministratorID uint) (*dtos.ProductResponse, error) {
	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	product, err := pu.ProductRepository.FindOne(int(id))
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check AdministratorID
	if product.AdministratorID != AdministratorID {
		return nil, errors.New("not allowed to update product")
	}

	if product.Name != req.Name {
		slug := helpers.CreateSlug(req.Name)

		checkSlug, err := pu.ProductRepository.FindOneBySlug(slug)
		if err == nil {
			if checkSlug.Slug == slug {
				slug = slug + "-" + helpers.RandomString(5)
			}
		}

		req.Slug = slug
		product.Slug = slug
	}

	if req.Name != "" && req.Detail != "" && req.Price != 0 && req.Stock != 0 && req.Category != "" {
		product.Name = req.Name
		product.Detail = req.Detail
		product.Price = req.Price
		product.Stock = req.Stock
		product.Category = req.Category
	}

	if req.Image != "" {
		product.Image = req.Image
	}

	updatedProduct, err := pu.ProductRepository.UpdateOne(product)
	if err != nil {
		return nil, errors.New("error updating product")
	}

	res := &dtos.ProductResponse{
		ID:       updatedProduct.ID,
		Slug:     updatedProduct.Slug,
		Name:     updatedProduct.Name,
		Detail:   updatedProduct.Detail,
		Price:    updatedProduct.Price,
		Stock:    updatedProduct.Stock,
		Category: updatedProduct.Category,
		Image:    updatedProduct.Image,
	}

	return res, nil
}

// DeleteProduct godoc
// @Summary      Delete Product
// @Description  Delete Product
// @Tags         Admin - Product
// @Accept       json
// @Produce      json
// @Param id query int false "Product ID"
// @Param        request body dtos.UpdateProductRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.ProductDeletedStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/product/{id} [delete]
// @Security BearerAuth
func (pu *productUsecase) DeleteOne(ctx context.Context, id, AdministratorID uint) error {
	_, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	product, err := pu.ProductRepository.FindOne(int(id))
	if err != nil {
		return errors.New("product not found")
	}

	if product.AdministratorID != AdministratorID {
		return errors.New("not allowed to delete product")
	}

	err = pu.ProductRepository.DeleteOne(product)
	if err != nil {
		return errors.New("error deleting product")
	}

	return nil
}
