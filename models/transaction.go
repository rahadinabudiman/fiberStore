package models

import (
	"context"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID     uint `json:"user_id" form:"user_id"`
	User       User `gorm:"foreignKey:UserID" json:"-"`
	CartID     uint `json:"cart_id" form:"cart_id"`
	Cart       Cart `gorm:"foreignKey:CartID" json:"-"`
	TotalPrice int  `json:"total_price" form:"total_price"`
}

type TransactionRepository interface {
	InsertOne(req *Transaction) (*Transaction, error)
	FindOne(cartID uint) (res *Transaction, err error)
	FindAll(userID uint) (*[]Transaction, error)
	BeginTx(ctx context.Context) *gorm.DB
	UpdateOne(req *Transaction, id uint) (res *Transaction, err error)
	DeleteOne(req *Transaction) error
}

type TransactionUsecase interface {
	InsertOne(ctx context.Context, req *Transaction) (*Transaction, error)
	FindOne(ctx context.Context, cartID uint) (*Transaction, error)
	FindAll(ctx context.Context, userID uint) (*[]Transaction, error)
}
