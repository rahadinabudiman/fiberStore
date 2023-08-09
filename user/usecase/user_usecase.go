package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/helpers"
	"fiberStore/models"
	"fiberStore/user/repository"
	"sort"
	"strings"
	"time"
)

type UserUsecase interface {
	InsertOne(ctx context.Context, req *dtos.UserRegister) (res *dtos.UserRegisterResponse, err error)
	FindOneById(ctx context.Context, id int) (res *dtos.UserProfileResponse, err error)
	FindAll(ctx context.Context, page, limit int, search, sortBy string) (*[]dtos.UserDetailResponse, int, error)
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

func (uu *userUsecase) FindOneById(ctx context.Context, id int) (res *dtos.UserProfileResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	Profile, err := uu.UserRepository.FindOneById(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	res = &dtos.UserProfileResponse{
		Name:     Profile.Name,
		Username: Profile.Username,
	}

	return res, nil
}

func (uu *userUsecase) FindAll(ctx context.Context, page, limit int, search, sortBy string) (res *[]dtos.UserDetailResponse, count int, err error) {
	var userResponses []dtos.UserDetailResponse

	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	users, count, err := uu.UserRepository.FindAll(page, limit, search)
	if err != nil {
		return res, 0, errors.New("error getting users")
	}

	for _, user := range *users {
		userResponse := dtos.UserDetailResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
		}
		userResponses = append(userResponses, userResponse)
	}

	if strings.ToLower(sortBy) == "asc" {
		sort.SliceStable(userResponses, func(i, j int) bool {
			return userResponses[i].Name < userResponses[j].Name
		})
	} else if strings.ToLower(sortBy) == "desc" {
		sort.SliceStable(userResponses, func(i, j int) bool {
			return userResponses[i].Name > userResponses[j].Name
		})
	}

	res = &userResponses

	return res, count, nil
}

func (uu *userUsecase) UpdateOne(req *models.User) (*models.User, error) {
	return uu.UserRepository.UpdateOne(req)
}

func (uu *userUsecase) DeleteOne(id int) error {
	return uu.UserRepository.DeleteOne(id)
}
