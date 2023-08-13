package models

import (
	"context"
	"fiberStore/dtos"

	"gorm.io/gorm"
)

type CartDetail struct {
	gorm.Model
	CartID    uint    `json:"cart_id" form:"cart_id"`
	Cart      Cart    `gorm:"foreignKey:CartID" json:"-"`
	UserID    uint    `json:"user_id" form:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"-"`
	ProductID uint    `json:"product_id" form:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
	Quantity  int     `json:"quantity" form:"quantity"`
}

type CartDetailRepository interface {
	InsertOne(req *CartDetail) (*CartDetail, error)
	FindOne(userID, productID uint) (res *CartDetail, err error)
	FindAll(userID uint) (*[]CartDetail, error)
	UpdateOne(req *CartDetail, id uint) (res *CartDetail, err error)
	DeleteOne(req *CartDetail) error
}

type CartDetailUsecase interface {
	InsertOne(ctx context.Context, req *dtos.InsertCartDetailRequest) (*dtos.InsertCartDetailResponse, error)
	FindAll(ctx context.Context, userID uint) (*dtos.CartDetailResponse, error)
	DeleteProduct(ctx context.Context, userID uint, productID uint) error
}
