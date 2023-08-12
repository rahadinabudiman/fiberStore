package models

import (
	"context"
	"fiberStore/dtos"

	"gorm.io/gorm"
)

type UserAmount struct {
	gorm.Model
	UserID uint    `gorm:"foreignKey:UserAmountRefer" json:"user_id" form:"user_id"`
	User   User    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Amount float64 `json:"amount" form:"amount"`
}

type UserAmountRepository interface {
	InsertOne(req *UserAmount) (*UserAmount, error)
	FindOne(id uint) (res *UserAmount, err error)
	UpdateOne(req *UserAmount, id uint) (res *UserAmount, err error)
	DeleteOne(req *UserAmount) error
}

type UserAmountUsecase interface {
	TopUpSaldo(ctx context.Context, req *dtos.TopUpSaldoRequest) (res *dtos.TopUpSaldoResponse, err error)
}
