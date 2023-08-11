package models

import "gorm.io/gorm"

type UserAmount struct {
	gorm.Model
	UserID uint    `gorm:"not null" json:"user_id" form:"user_id"`
	Amount float64 `gorm:"not null" json:"amount" form:"amount"`
}
