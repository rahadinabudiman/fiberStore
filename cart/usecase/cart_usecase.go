package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/models"
	"time"
)

type cartUsecase struct {
	CartRepository    models.CartRepository
	ProductRepository models.ProductRepository
	UserRepository    models.UserRepository
	contextTimeout    time.Duration
}

func NewCartUsecase(CartRepository models.CartRepository, ProductRepository models.ProductRepository, UserRepository models.UserRepository, contextTimeout time.Duration) models.CartUsecase {
	return &cartUsecase{
		CartRepository:    CartRepository,
		ProductRepository: ProductRepository,
		UserRepository:    UserRepository,
		contextTimeout:    contextTimeout,
	}
}

func (cu *cartUsecase) InsertOne(ctx context.Context, req *dtos.InsertCartRequest) (*dtos.InsertCartResponse, error) {
	_, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Check User Exist
	user, err := cu.UserRepository.FindOneById(int(req.UserID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check Product Exist
	product, err := cu.ProductRepository.FindOne(int(req.ProductID))
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check Product Stock
	if product.Stock == 0 {
		return nil, errors.New("product stock empty")
	}

	if product.Stock < int64(req.Quantity) {
		return nil, errors.New("product stock not enough")
	}

	// Check Cart Exist
	cart, _ := cu.CartRepository.FindOne(req.UserID, req.ProductID)

	if cart != nil {
		// Update Cart
		cart.Quantity += req.Quantity

		if product.Stock < int64(cart.Quantity) {
			return nil, errors.New("product stock not enough")
		}

		_, err = cu.CartRepository.UpdateOne(cart, cart.ID)
		if err != nil {
			return nil, errors.New("failed to update cart")
		}

		res := &dtos.InsertCartResponse{
			Name:        user.Name,
			ProductName: product.Name,
			Quantity:    req.Quantity,
		}

		return res, nil

	} else {
		// Insert Cart
		cart := &models.Cart{
			UserID:    req.UserID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		_, err = cu.CartRepository.InsertOne(cart)
		if err != nil {
			return nil, errors.New("failed to insert cart")
		}

		res := &dtos.InsertCartResponse{
			Name:        user.Name,
			ProductName: product.Name,
			Quantity:    req.Quantity,
		}

		return res, nil
	}
}

func (cu *cartUsecase) FindAll(ctx context.Context, userID uint) (*dtos.CartResponse, error) {
	var detailCartResponses []dtos.DetailCartResponse

	_, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	carts, err := cu.CartRepository.FindAll(userID)
	if err != nil {
		return nil, errors.New("failed to get carts")
	}

	grandTotal := 0

	if len(*carts) == 0 {
		return nil, errors.New("cart is empty, please add product to cart first")
	} else {

		for _, cart := range *carts {
			product, err := cu.ProductRepository.FindOne(int(cart.ProductID))
			if err != nil {
				return nil, errors.New("product not found")
			}

			totalPrice := int(product.Price) * cart.Quantity
			grandTotal += totalPrice

			cartResponse := dtos.DetailCartResponse{
				ID:          product.ID,
				ProductName: product.Name,
				Price:       int(product.Price),
				Quantity:    cart.Quantity,
				TotalPrice:  int(product.Price) * cart.Quantity,
			}
			detailCartResponses = append(detailCartResponses, cartResponse)
		}

		cartResponses := dtos.CartResponse{
			Produks:    detailCartResponses,
			GrandTotal: grandTotal,
		}
		return &cartResponses, nil
	}
}

func (cu *cartUsecase) DeleteProduct(ctx context.Context, userID uint, productID uint) error {
	_, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	_, err := cu.UserRepository.FindOneById(int(userID))
	if err != nil {
		return errors.New("user not found")
	}

	_, err = cu.ProductRepository.FindOne(int(productID))
	if err != nil {
		return errors.New("product not found")
	}

	product, err := cu.CartRepository.FindAll(userID)
	if err != nil {
		return errors.New("please add product to cart first")
	}

	for _, cart := range *product {
		if cart.ProductID != productID {
			return errors.New("product not found in cart")
		}
	}

	cart, err := cu.CartRepository.FindOne(userID, productID)
	if err != nil {
		return errors.New("please add product to cart first")
	}

	err = cu.CartRepository.DeleteOne(cart)
	if err != nil {
		return errors.New("failed to delete product in cart")
	}

	return nil
}
