package dtos

type InsertCartDetailRequest struct {
	CartID    uint  `json:"cart_id" form:"cart_id"`
	UserID    uint  `json:"user_id" form:"user_id"`
	ProductID uint  `json:"product_id" form:"product_id"`
	Quantity  int64 `json:"quantity" form:"quantity"`
}

type InsertCartDetailResponse struct {
	CartID      uint   `json:"cart_id" form:"cart_id"`
	Name        string `json:"name" form:"name"`
	ProductName string `json:"product_name" form:"product_name"`
	Quantity    int64  `json:"quantity" form:"quantity"`
}

type AddProductToCart struct {
	ProductID uint `json:"product_id" form:"product_id"`
	Quantity  int  `json:"quantity" form:"quantity"`
}

type DetailCartDetailResponse struct {
	ID          uint   `json:"id" form:"id"`
	ProductName string `json:"product_name" form:"product_name"`
	Price       int64  `json:"price" form:"price"`
	Quantity    int64  `json:"quantity" form:"quantity"`
	TotalPrice  int64  `json:"total_price" form:"total_price"`
}

type CartDetailResponse struct {
	Produks    []DetailCartDetailResponse `json:"produks"`
	GrandTotal int                        `json:"grand_total"`
}
