package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/helpers"
	"fiberStore/models"
	"fiberStore/user/repository"
	"time"
)

type UserUsecase interface {
	InsertOne(ctx context.Context, req *dtos.UserRegister) (res *dtos.UserRegisterResponse, err error)
	FindOneById(id int) (*models.User, error)
	FindAll(page, limit int, search string) (*[]models.User, int, error)
	UpdateOne(req *models.User) (*models.User, error)
	DeleteOne(id int) error
}

type userUsecase struct {
	UserRepository repository.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(UserRepository repository.UserRepository, contextTimeout time.Duration) UserUsecase {
	return &userUsecase{
		UserRepository: UserRepository,
		contextTimeout: contextTimeout,
	}
}

func (uu *userUsecase) InsertOne(ctx context.Context, req *dtos.UserRegister) (res *dtos.UserRegisterResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	_, err = uu.UserRepository.FindOneByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password and confirm password not match")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	CreateUser := &models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: passwordHash,
	}

	resp, err := uu.UserRepository.InsertOne(CreateUser)
	if err != nil {
		return nil, errors.New("error creating user")
	}

	res = &dtos.UserRegisterResponse{
		Name:      resp.Name,
		Username:  resp.Username,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}

	return res, nil
}

func (uu *userUsecase) FindOneById(id int) (*models.User, error) {
	return uu.UserRepository.FindOneById(id)
}

func (uu *userUsecase) FindAll(page, limit int, search string) (*[]models.User, int, error) {
	return uu.UserRepository.FindAll(page, limit, search)
}

func (uu *userUsecase) UpdateOne(req *models.User) (*models.User, error) {
	return uu.UserRepository.UpdateOne(req)
}

func (uu *userUsecase) DeleteOne(id int) error {
	return uu.UserRepository.DeleteOne(id)
}
