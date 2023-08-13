package dtos

type InsertTransactionRequest struct {
	CartID uint `json:"cart_id" form:"cart_id"`
}
