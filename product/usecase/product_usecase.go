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

func (pu *productUsecase) FindAll(ctx context.Context, page, limit int) (*[]dtos.ProductResponse, int, error) {
	return nil, 0, nil
}

func (pu *productUsecase) FindQueryAll(ctx context.Context, page, limit int, search string) (*[]dtos.ProductResponse, int, error) {
	return nil, 0, nil
}

func (pu *productUsecase) FindOne(ctx context.Context, id uint) (*dtos.ProductResponse, error) {
	return nil, nil
}

func (pu *productUsecase) UpdateOne(ctx context.Context, req *dtos.ProductRequest, id uint) (*dtos.ProductResponse, error) {
	return nil, nil
}

func (pu *productUsecase) DeleteOne(ctx context.Context, id uint) error {
	return nil
}
