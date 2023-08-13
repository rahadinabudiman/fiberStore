package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	CartID       uint       `json:"cart_id" form:"cart_id"`
	Cart         Cart       `gorm:"foreignKey:CartID" json:"-"`
	CartDetailID uint       `json:"cart_detail_id" form:"cart_detail_id"`
	CartDetail   CartDetail `gorm:"foreignKey:CartDetailID" json:"-"`
	TotalPrice   int        `json:"total_price" form:"total_price"`
}
