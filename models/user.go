package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Role     string `gorm:"type:enum('Customer', 'Admin');default:'Customer'; not-null" example:"Admin" json:"role" form:"role"`
}
