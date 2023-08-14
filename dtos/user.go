package dtos

import "time"

type UserRegister struct {
	Name            string `json:"name" form:"name" binding:"required"`
	Username        string `gorm:"unique" json:"username" form:"username" binding:"required"`
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required"`
	Role            string `gorm:"type:enum('Customer', 'Admin');default:'Customer'; not-null" example:"Admin" json:"role" form:"role"`
}

type UserUpdateProfile struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
}

type DeleteUserRequest struct {
	Password string `json:"password" form:"password" validate:"gte=6" example:"rahadinabudimansundara"`
}

type UserLoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type TopUpSaldoRequest struct {
	Username string  `json:"username" form:"username"`
	Amount   float64 `json:"amount" form:"amount" validate:"required" example:"100000"`
}

type TopUpSaldoResponse struct {
	Name   string  `json:"name" form:"name"`
	Amount float64 `json:"amount" form:"amount"`
}

type UserLoginResponse struct {
	Username string `json:"username" form:"username" binding:"required" example:"rahadinabudimansundara"`
	Token    string `json:"token" form:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

type UserRegisterResponse struct {
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserProfileResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserDetailResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Status   string `json:"status"`
}
