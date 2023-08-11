package models

import "gorm.io/gorm"

type UserAmount struct {
	gorm.Model
	UserID uint    `json:"user_id" form:"user_id"`
	Amount float64 `json:"amount" form:"amount"`
}