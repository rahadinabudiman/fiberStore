package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/models"
	"time"
)

type userAmountUsecase struct {
	UserAmountRepository models.UserAmountRepository
	UserRepository       models.UserRepository
	contextTimeout       time.Duration
}

func NewUserAmountUsecase(UserAmountRepository models.UserAmountRepository, UserRepository models.UserRepository, contextTimeout time.Duration) models.UserAmountUsecase {
	return &userAmountUsecase{
		UserAmountRepository: UserAmountRepository,
		UserRepository:       UserRepository,
		contextTimeout:       contextTimeout,
	}
}

func (uau *userAmountUsecase) TopUpSaldo(ctx context.Context, req *dtos.TopUpSaldoRequest) (res *dtos.TopUpSaldoResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uau.contextTimeout)
	defer cancel()

	user, err := uau.UserRepository.FindOneByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username not found")
	}

	userAmount, err := uau.UserAmountRepository.FindOne(user.ID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	userAmount.Amount += req.Amount

	_, err = uau.UserAmountRepository.UpdateOne(userAmount, userAmount.UserID)
	if err != nil {
		return nil, errors.New("failed to top up saldo")
	}

	res = &dtos.TopUpSaldoResponse{
		Name:   user.Name,
		Amount: req.Amount,
	}

	return res, nil
}
