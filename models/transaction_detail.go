package models

import (
	"context"
	"fiberStore/dtos"

	"gorm.io/gorm"
)

type TransactionDetail struct {
	gorm.Model
	TransactionID uint        `json:"transaction_id" form:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"-"`
	UserID        uint        `json:"user_id" form:"user_id"`
	User          User        `gorm:"foreignKey:UserID" json:"-"`
	CartID        uint        `json:"cart_id" form:"cart_id"`
	Cart          Cart        `gorm:"foreignKey:CartID" json:"-"`
	ProductID     uint        `json:"product_id" form:"product_id"`
	Product       Product     `gorm:"foreignKey:ProductID" json:"-"`
	Quantity      int64       `json:"quantity" form:"quantity"`
	TotalPrice    int64       `json:"total_price" form:"total_price"`
}

type TransactionDetailRepository interface {
	InsertOne(req *TransactionDetail) (*TransactionDetail, error)
	FindOne(cartID uint) (res *TransactionDetail, err error)
	FindAll(userID uint) (*[]TransactionDetail, error)
	BeginTx(ctx context.Context) *gorm.DB
	UpdateOne(req *TransactionDetail, id uint) (res *TransactionDetail, err error)
	DeleteOne(req *TransactionDetail) error
}

type TransactionDetailUsecase interface {
	InsertOne(ctx context.Context, req *TransactionDetail) (*dtos.InsertTransactionDetailResponse, error)
	FindOne(ctx context.Context, cartID uint) (*TransactionDetail, error)
	FindAll(ctx context.Context, userID uint) (*[]TransactionDetail, error)
}
