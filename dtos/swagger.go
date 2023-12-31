package dtos

import "fiberStore/utils"

// Authentikasi
type LoginStatusOKResponse struct {
	StatusCode int               `json:"status_code" example:"200"`
	Message    string            `json:"message" form:"message" example:"Login Success"`
	Data       UserLoginResponse `json:"data"`
}

type RegisterStatusOKResponse struct {
	StatusCode int                  `json:"status_code" example:"200"`
	Message    string               `json:"message" form:"message" example:"Register Success"`
	Data       UserRegisterResponse `json:"data"`
}

// Profile User
type ProfileStatusOKResponse struct {
	StatusCode int                 `json:"status_code" example:"200"`
	Message    string              `json:"message" form:"message" example:"Get Profile Success"`
	Data       UserProfileResponse `json:"data"`
}

type GetAllUserStatusOKResponse struct {
	StatusCode int                `json:"status_code" example:"200"`
	Message    string             `json:"message" example:"Successfully get profile"`
	Data       UserDetailResponse `json:"data"`
	Meta       utils.Meta         `json:"meta"`
}

// Product
type InsertProductStatusOKResponse struct {
	StatusCode int                  `json:"status_code" example:"200"`
	Message    string               `json:"message" example:"Successfully insert product"`
	Data       InserProductResponse `json:"data"`
}

type ProductStatusOKResponse struct {
	StatusCode int                  `json:"status_code" example:"200"`
	Message    string               `json:"message" example:"Successfully insert product"`
	Data       InserProductResponse `json:"data"`
}

type ProductDeletedStatusOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" example:"Successfully deleted product"`
}

type GetAllProductStatusOKResponse struct {
	StatusCode int               `json:"status_code" example:"200"`
	Message    string            `json:"message" example:"Successfully get all product"`
	Data       []ProductResponse `json:"data"`
	Meta       utils.Meta        `json:"meta"`
}

// Cart
type InsertCartStatusOKResponse struct {
	StatusCode int              `json:"status_code" example:"200"`
	Message    string           `json:"message" example:"Successfully insert cart"`
	Data       AddProductToCart `json:"data"`
}

type CartStatusOKResponse struct {
	StatusCode int                `json:"status_code" example:"200"`
	Message    string             `json:"message" example:"Successfully get cart"`
	Data       CartDetailResponse `json:"data"`
}

type CartDeletedStatusOKResponse struct {
	StatusCode int    `json:"status_code" example:"200"`
	Message    string `json:"message" example:"Successfully deleted cart"`
}

// Transaction
type InsertTransactionStatusOKResponse struct {
	StatusCode int                             `json:"status_code" example:"200"`
	Message    string                          `json:"message" example:"Successfully insert transaction"`
	Data       InsertTransactionDetailResponse `json:"data"`
}

type TransactionStatusOKResponse struct {
	StatusCode int                             `json:"status_code" example:"200"`
	Message    string                          `json:"message" example:"Successfully get transaction"`
	Data       DetailTransactionDetailResponse `json:"data"`
}

// Top Up Saldo
type TopUpStatusOKResponse struct {
	StatusCode int                `json:"status_code" example:"200"`
	Message    string             `json:"message" example:"Successfully top up"`
	Data       TopUpSaldoResponse `json:"data"`
}

type StatusOKResponse struct {
	StatusCode int         `json:"status_code" example:"200"`
	Message    string      `json:"message" example:"Successfully"`
	Data       interface{} `json:"data"`
}

type StatusOKDeletedResponse struct {
	StatusCode int         `json:"status_code" example:"200"`
	Message    string      `json:"message" example:"Successfully deleted"`
	Errors     interface{} `json:"errors"`
}

type BadRequestResponse struct {
	StatusCode int         `json:"status_code" example:"400"`
	Message    string      `json:"message" example:"Bad Request"`
	Errors     interface{} `json:"errors"`
}

type UnauthorizedResponse struct {
	StatusCode int         `json:"status_code" example:"401"`
	Message    string      `json:"message" example:"Unauthorized"`
	Errors     interface{} `json:"errors"`
}

type ForbiddenResponse struct {
	StatusCode int         `json:"status_code" example:"403"`
	Message    string      `json:"message" example:"Forbidden"`
	Errors     interface{} `json:"errors"`
}

type NotFoundResponse struct {
	StatusCode int         `json:"status_code" example:"404"`
	Message    string      `json:"message" example:"Not Found"`
	Errors     interface{} `json:"errors"`
}

type InternalServerErrorResponse struct {
	StatusCode int         `json:"status_code" example:"500"`
	Message    string      `json:"message" example:"Internal Server Error"`
	Errors     interface{} `json:"errors"`
}
