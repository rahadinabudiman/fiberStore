package models

import (
	"context"
	"fiberStore/dtos"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                  string    `json:"name" form:"name" binding:"required"`
	Username              string    `json:"username" form:"username" binding:"required"`
	Password              string    `json:"password" form:"password" binding:"required"`
	Role                  string    `gorm:"type:enum('Customer', 'Admin');default:'Customer'; not-null" example:"Admin" json:"role" form:"role"`
	AdministratorProducts []Product `gorm:"foreignKey:AdministratorID" json:"-"`
}

type UserRepository interface {
	InsertOne(req *User) (*User, error)
	FindOneById(id int) (*User, error)
	FindOneAdmin(id uint) (*User, error)
	FindOneByUsername(username string) (*User, error)
	FindAll(page, limit int, search string) (*[]User, int, error)
	UpdateOne(req *User) (*User, error)
	DeleteOne(user *User) error
}

type UserUsecase interface {
	// Authentikasi User
	LoginUser(ctx context.Context, c *fiber.Ctx, req *dtos.UserLoginRequest) (res *dtos.UserLoginResponse, err error)
	// CRUD User
	InsertOne(ctx context.Context, req *dtos.UserRegister) (res *dtos.UserRegisterResponse, err error)
	FindOneById(ctx context.Context, id int) (res *dtos.UserProfileResponse, err error)
	FindAll(ctx context.Context, page, limit int, search, sortBy string) (*[]dtos.UserDetailResponse, int, error)
	UpdateOne(ctx context.Context, id int, req *dtos.UserUpdateProfile) (res *dtos.UserProfileResponse, err error)
	UpdatePassword(ctx context.Context, id uint, req *dtos.UserUpdatePassword) (res *dtos.UserProfileResponse, err error)
	DeleteOne(ctx context.Context, id uint, req dtos.DeleteUserRequest) error
}
