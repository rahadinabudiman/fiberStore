package dtos

type InsertTransactionDetailRequest struct {
	CartID uint `json:"cart_id" form:"cart_id"`
}

type InsertTransactionDetailResponse struct {
	TransactionID uint                         `json:"transaction_id" form:"transaction_id"`
	Name          string                       `json:"name" form:"name"`
	Produk        []ProductTransactionResponse `json:"produk" form:"produk"`
	TotalPrice    int64                        `json:"total_price" form:"total_price"`
}
