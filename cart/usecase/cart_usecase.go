package usecase

import (
	"context"
	"errors"
	"fiberStore/models"
	"time"
)

type cartUsecase struct {
	CartRepository models.CartRepository
	contextTimeout time.Duration
}

func NewCartUsecase(CartRepository models.CartRepository, contextTimeout time.Duration) models.CartUsecase {
	return &cartUsecase{
		CartRepository: CartRepository,
		contextTimeout: contextTimeout,
	}
}

func (cu *cartUsecase) InsertOne(ctx context.Context, userID uint) (*models.Cart, error) {
	_, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	req := &models.Cart{
		UserID: userID,
	}

	res, err := cu.CartRepository.InsertOne(req)
	if err != nil {
		return nil, errors.New("failed to add cart")
	}

	return res, nil
}

func (cu *cartUsecase) FindOne(ctx context.Context, userID uint) (*models.Cart, error) {
	_, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	res, err := cu.CartRepository.FindOne(userID)
	if err != nil {
		return nil, errors.New("cart not found, please add cart first")
	}

	return res, nil
}
