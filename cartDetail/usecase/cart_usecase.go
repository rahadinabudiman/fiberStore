package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/models"
	"time"
)

type CartDetailUsecase struct {
	CartDetailRepository models.CartDetailRepository
	CartRepository       models.CartRepository
	ProductRepository    models.ProductRepository
	UserRepository       models.UserRepository
	contextTimeout       time.Duration
}

func NewCartDetailUsecase(CartDetailRepository models.CartDetailRepository, CartRepository models.CartRepository, ProductRepository models.ProductRepository, UserRepository models.UserRepository, contextTimeout time.Duration) models.CartDetailUsecase {
	return &CartDetailUsecase{
		CartDetailRepository: CartDetailRepository,
		CartRepository:       CartRepository,
		ProductRepository:    ProductRepository,
		UserRepository:       UserRepository,
		contextTimeout:       contextTimeout,
	}
}

// AddProductToCart godoc
// @Summary      Add Product To Cart
// @Description  Add Product To Cart
// @Tags         User - Cart
// @Accept       json
// @Produce      json
// @Param        request body dtos.AddProductToCart true "Payload Body [RAW]"
// @Success      200 {object} dtos.InsertCartStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /cart [post]
// @Security BearerAuth
func (cu *CartDetailUsecase) InsertOne(ctx context.Context, req *dtos.InsertCartDetailRequest) (*dtos.InsertCartDetailResponse, error) {
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

	// Check Cart Exist or Create New One
	cart, err := cu.CartRepository.FindOne(req.UserID)
	if err != nil {
		cart = &models.Cart{UserID: req.UserID}
		newCart, err := cu.CartRepository.InsertOne(cart)
		if err != nil {
			return nil, errors.New("failed to insert Cart")
		}
		cart = newCart
	}

	// Check Product Stock
	if product.Stock == 0 {
		return nil, errors.New("product stock empty")
	}

	if product.Stock < int64(req.Quantity) {
		return nil, errors.New("product stock not enough")
	}

	// Check CartDetail Exist
	CartDetail, _ := cu.CartDetailRepository.FindOne(req.UserID, req.ProductID)
	if CartDetail != nil {
		// Update CartDetail
		CartDetail.Quantity += req.Quantity

		if product.Stock < int64(CartDetail.Quantity) {
			return nil, errors.New("product stock not enough")
		}

		_, err = cu.CartDetailRepository.UpdateOne(CartDetail, CartDetail.ID)
		if err != nil {
			return nil, errors.New("failed to update CartDetail")
		}

		res := &dtos.InsertCartDetailResponse{
			CartID:      cart.ID,
			Name:        user.Name,
			ProductName: product.Name,
			Quantity:    req.Quantity,
		}

		return res, nil

	} else {
		// Insert CartDetail
		CartDetail := &models.CartDetail{
			CartID:    cart.ID,
			UserID:    req.UserID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		_, err = cu.CartDetailRepository.InsertOne(CartDetail)
		if err != nil {
			return nil, errors.New("failed to insert CartDetail")
		}

		res := &dtos.InsertCartDetailResponse{
			CartID:      cart.ID,
			Name:        user.Name,
			ProductName: product.Name,
			Quantity:    req.Quantity,
		}

		return res, nil
	}
}

// GetProductInCart godoc
// @Summary      Get Product In Cart
// @Description  Get Product In Cart
// @Tags         User - Cart
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.InsertCartStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /cart [get]
// @Security BearerAuth
func (cu *CartDetailUsecase) FindAll(ctx context.Context, userID uint) (*dtos.CartDetailResponse, error) {
	var detailCartDetailResponses []dtos.DetailCartDetailResponse

	_, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	CartDetails, err := cu.CartDetailRepository.FindAll(userID)
	if err != nil {
		return nil, errors.New("failed to get CartDetails")
	}

	grandTotal := 0

	if len(*CartDetails) == 0 {
		return nil, errors.New("CartDetail is empty, please add product to CartDetail first")
	} else {

		for _, CartDetail := range *CartDetails {
			product, err := cu.ProductRepository.FindOne(int(CartDetail.ProductID))
			if err != nil {
				return nil, errors.New("product not found")
			}

			totalPrice := product.Price * CartDetail.Quantity
			grandTotal += int(totalPrice)

			CartDetailResponse := dtos.DetailCartDetailResponse{
				ID:          product.ID,
				ProductName: product.Name,
				Price:       product.Price,
				Quantity:    CartDetail.Quantity,
				TotalPrice:  product.Price * CartDetail.Quantity,
			}
			detailCartDetailResponses = append(detailCartDetailResponses, CartDetailResponse)
		}

		CartDetailResponses := dtos.CartDetailResponse{
			Produks:    detailCartDetailResponses,
			GrandTotal: grandTotal,
		}
		return &CartDetailResponses, nil
	}
}

// DeleteProductInCart godoc
// @Summary      Delete Product In Cart
// @Description  Delete Product In Cart
// @Tags         User - Cart
// @Accept       json
// @Produce      json
// @Param product_id query int false "Product number"
// @Success      200 {object} dtos.CartDeletedStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /cart [delete]
// @Security BearerAuth
func (cu *CartDetailUsecase) DeleteProduct(ctx context.Context, userID uint, productID uint) error {
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

	product, err := cu.CartDetailRepository.FindAll(userID)
	if err != nil {
		return errors.New("please add product to CartDetail first")
	}

	for _, CartDetail := range *product {
		if CartDetail.ProductID != productID {
			return errors.New("product not found in CartDetail")
		}
	}

	CartDetail, err := cu.CartDetailRepository.FindOne(userID, productID)
	if err != nil {
		return errors.New("please add product to CartDetail first")
	}

	err = cu.CartDetailRepository.DeleteOne(CartDetail)
	if err != nil {
		return errors.New("failed to delete product in CartDetail")
	}

	return nil
}
