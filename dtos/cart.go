package dtos

type InsertCartDetailRequest struct {
	CartID    uint `json:"cart_id" form:"cart_id"`
	UserID    uint `json:"user_id" form:"user_id"`
	ProductID uint `json:"product_id" form:"product_id"`
	Quantity  int  `json:"quantity" form:"quantity"`
}

type InsertCartDetailResponse struct {
	CartID      uint   `json:"cart_id" form:"cart_id"`
	Name        string `json:"name" form:"name"`
	ProductName string `json:"product_name" form:"product_name"`
	Quantity    int    `json:"quantity" form:"quantity"`
}

type DetailCartDetailResponse struct {
	ID          uint   `json:"id" form:"id"`
	ProductName string `json:"product_name" form:"product_name"`
	Price       int    `json:"price" form:"price"`
	Quantity    int    `json:"quantity" form:"quantity"`
	TotalPrice  int    `json:"total_price" form:"total_price"`
}

type CartDetailResponse struct {
	Produks    []DetailCartDetailResponse `json:"produks"`
	GrandTotal int                        `json:"grand_total"`
}
