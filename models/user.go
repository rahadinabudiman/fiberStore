package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `gorm:"unique" json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Role     string `gorm:"type:ENUM('customer','admin')" default:"customer" json:"role" form:"role"`
}
