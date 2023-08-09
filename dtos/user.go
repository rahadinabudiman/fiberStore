package dtos

import "time"

type UserRegister struct {
	Name            string `json:"name" form:"name" binding:"required"`
	Username        string `gorm:"unique" json:"username" form:"username" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required"`
	Role            string `gorm:"type:enum('Customer', 'Admin');default:'Customer'; not-null" example:"Admin" json:"role" form:"role"`
}

type UserRegisterResponse struct {
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
