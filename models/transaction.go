package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID uint `json:"user_id" form:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`
}

type TransactionRepository interface {
	InsertOne(req *Transaction) (*Transaction, error)
	FindOne(userID uint) (res *Transaction, err error)
	DeleteOne(req *Transaction) error
}
