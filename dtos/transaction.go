package dtos

import "time"

type InsertTransactionDetailRequest struct {
	CartID uint `json:"cart_id" form:"cart_id"`
}

type InsertTransactionDetailResponse struct {
	TransactionID uint                         `json:"transaction_id" form:"transaction_id"`
	Name          string                       `json:"name" form:"name"`
	Produk        []ProductTransactionResponse `json:"produk" form:"produk"`
	TotalPrice    int64                        `json:"total_price" form:"total_price"`
}

type DetailTransactionDetailResponse struct {
	ID          uint      `json:"id" form:"id"`
	CreatedAt   time.Time `json:"created_at" form:"created_at"`
	ProductName string    `json:"product_name" form:"product_name"`
	Price       int64     `json:"price" form:"price"`
	Quantity    int64     `json:"quantity" form:"quantity"`
	TotalPrice  int64     `json:"total_price" form:"total_price"`
}

type TransactionDetailRequest struct {
	UserID uint `json:"user_id" form:"user_id"`
}
