package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/helpers"
	"fiberStore/middlewares"
	"fiberStore/models"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	UserRepository       models.UserRepository
	UserAmountRepository models.UserAmountRepository
	contextTimeout       time.Duration
}

func NewUserUsecase(UserRepository models.UserRepository, UserAmountRepository models.UserAmountRepository, contextTimeout time.Duration) models.UserUsecase {
	return &userUsecase{
		UserRepository:       UserRepository,
		UserAmountRepository: UserAmountRepository,
		contextTimeout:       contextTimeout,
	}
}

// UserLogin godoc
// @Summary      Login User with Username and Password
// @Description  Login an account
// @Tags         Authentikasi
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserLoginRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.LoginStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /login [post]
func (uu *userUsecase) LoginUser(ctx context.Context, c *fiber.Ctx, req *dtos.UserLoginRequest) (res *dtos.UserLoginResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.UserRepository.FindOneByUsername(req.Username)
	if err != nil {
		return nil, errors.New("username not found")
	}

	if !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("account has been deleted")
	}

	fmt.Println(user)

	fmt.Println(user.DeletedAt.Time.IsZero())

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil {
		return nil, errors.New("username or password is wrong")
	}

	token, err := middlewares.CreateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	middlewares.CreateCookieToken(c, token)

	res = &dtos.UserLoginResponse{
		Username: req.Username,
		Token:    token,
	}
	return res, nil
}

// RegisterLogin godoc
// @Summary      Register User with Username and Password
// @Description  Register an account
// @Tags         Authentikasi
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserRegister true "Payload Body [RAW]"
// @Success      200 {object} dtos.RegisterStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /register [post]
func (uu *userUsecase) InsertOne(ctx context.Context, req *dtos.UserRegister) (res *dtos.UserRegisterResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	_, err = uu.UserRepository.FindOneByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	if strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.ConfirmPassword) == "" {
		return nil, errors.New("password and confirm password cannot be empty")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
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

	userAmount := &models.UserAmount{
		UserID: resp.ID,
		Amount: 0.0,
	}

	_, err = uu.UserAmountRepository.InsertOne(userAmount)
	if err != nil {
		return nil, errors.New("error creating user amount")
	}

	return res, nil
}

// GetProfile godoc
// @Summary      Get Profile
// @Description  User get profile
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.ProfileStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user [get]
// @Security BearerAuth
func (uu *userUsecase) FindOneById(ctx context.Context, id int) (res *dtos.UserProfileResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	Profile, err := uu.UserRepository.FindOneById(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !Profile.DeletedAt.Time.IsZero() {
		return nil, errors.New("account has been deleted")
	}

	res = &dtos.UserProfileResponse{
		Name:     Profile.Name,
		Username: Profile.Username,
	}

	return res, nil
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get all users
// @Tags         Admin - User
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search data"
// @Param sortBy query string false "Sort by name" Enums(asc, desc)
// @Success      200 {object} dtos.GetAllUserStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/user [get]
// @Security BearerAuth
func (uu *userUsecase) FindAll(ctx context.Context, page, limit int, search, sortBy string) (res *[]dtos.UserDetailResponse, count int, err error) {
	var userResponses []dtos.UserDetailResponse

	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	users, count, err := uu.UserRepository.FindAll(page, limit, search)
	if err != nil {
		return res, 0, errors.New("error getting users")
	}

	for _, user := range *users {
		Status := "Active"
		if !user.DeletedAt.Time.IsZero() {
			Status = "Deleted"
		}

		userResponse := dtos.UserDetailResponse{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Status:   Status,
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

// UserUpdateProfile godoc
// @Summary      Update Profile
// @Description  User update an Profile
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserUpdateProfile true "Payload Body [RAW]"
// @Success      200 {object} dtos.ProfileStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user [put]
// @Security BearerAuth
func (uu *userUsecase) UpdateOne(ctx context.Context, id int, req *dtos.UserUpdateProfile) (res *dtos.UserProfileResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.UserRepository.FindOneById(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Username != "" {
		user.Username = req.Username
	}

	user.Name = req.Name
	user.Username = req.Username

	user, err = uu.UserRepository.UpdateOne(user)
	if err != nil {
		return nil, errors.New("error updating user")
	}

	res = &dtos.UserProfileResponse{
		Name:     user.Name,
		Username: user.Username,
	}

	return res, nil
}

// USerChangePassword godoc
// @Summary      Change Password
// @Description  User change PAssword
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.UserUpdatePassword true "Payload Body [RAW]"
// @Success      200 {object} dtos.ProfileStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/password [put]
// @Security BearerAuth
func (uu *userUsecase) UpdatePassword(ctx context.Context, id uint, req *dtos.UserUpdatePassword) (res *dtos.UserProfileResponse, err error) {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.UserRepository.FindOneById(int(id))
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.OldPassword == "" || req.NewPassword == "" || req.ConfirmNewPassword == "" {
		return nil, errors.New("password cannot be empty")
	}

	if len(req.NewPassword) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	if req.NewPassword != req.ConfirmNewPassword {
		return nil, errors.New("password and confirm password not match")
	}

	err = helpers.ComparePassword(req.OldPassword, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("old password is incorrect")
	}

	passwordHash, err := helpers.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	user.Password = passwordHash

	user, err = uu.UserRepository.UpdateOne(user)
	if err != nil {
		return nil, errors.New("error updating user")
	}

	res = &dtos.UserProfileResponse{
		Name:     user.Name,
		Username: user.Username,
	}

	return res, nil
}

// UserDeleteProfile godoc
// @Summary      Delete Profile
// @Description  User delete an Profile
// @Tags         User - Account
// @Accept       json
// @Produce      json
// @Param        request body dtos.DeleteUserRequest true "Payload Body [RAW]"
// @Success      200 {object} dtos.StatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user [delete]
// @Security BearerAuth
func (uu *userUsecase) DeleteOne(ctx context.Context, id uint, req dtos.DeleteUserRequest) error {
	_, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.UserRepository.FindOneById(int(id))
	if err != nil {
		return errors.New("user not found")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	err = helpers.ComparePassword(req.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("password is incorrect")
	}

	err = uu.UserRepository.DeleteOne(user)
	if err != nil {
		return errors.New("error deleting user")
	}

	amount, err := uu.UserAmountRepository.FindOne(user.ID)
	if err != nil {
		return errors.New("error getting account balance")
	}

	err = uu.UserAmountRepository.DeleteOne(amount)
	if err != nil {
		return errors.New("error deleting account balance")
	}

	return nil
}
