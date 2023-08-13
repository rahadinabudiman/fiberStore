package models

import (
	"context"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID uint `json:"user_id" form:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`
}

type CartRepository interface {
	InsertOne(req *Cart) (*Cart, error)
	FindOne(userID uint) (res *Cart, err error)
	DeleteOne(req *Cart) error
}

type CartUsecase interface {
	InsertOne(ctx context.Context, userID uint) (*Cart, error)
	FindOne(ctx context.Context, userID uint) (*Cart, error)
}
