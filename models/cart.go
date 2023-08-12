package models

import (
	"context"
	"fiberStore/dtos"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserID    uint    `gorm:"foreignKey:CartRefer" json:"user_id" form:"user_id"`
	User      User    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ProductID uint    `gorm:"foreignKey:CartRefer" json:"product_id" form:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"-"`
	Quantity  int     `json:"quantity" form:"quantity"`
}

type CartRepository interface {
	InsertOne(req *Cart) (*Cart, error)
	FindOne(userID, productID uint) (res *Cart, err error)
	FindAll(userID uint) (*[]Cart, error)
	UpdateOne(req *Cart, id uint) (res *Cart, err error)
	DeleteOne(req *Cart) error
}

type CartUsecase interface {
	InsertOne(ctx context.Context, req *dtos.InsertCartRequest) (*dtos.InsertCartResponse, error)
	FindAll(ctx context.Context, userID uint) (*dtos.CartResponse, error)
	DeleteProduct(ctx context.Context, userID uint, productID uint) error
}
