package dtos

import "mime/multipart"

type ProductRequest struct {
	AdministratorID uint   `json:"administrator_id" form:"administrator_id"`
	Slug            string `json:"slug" form:"slug"`
	Name            string `json:"name" form:"name"`
	Detail          string `json:"detail" form:"detail"`
	Price           int64  `json:"price" form:"price"`
	Stock           int64  `json:"stock" form:"stock"`
	Category        string `json:"category" form:"category"`
	Image           string `json:"image" form:"image"`
}

type UpdateProductRequest struct {
	Name     string `json:"name" form:"name"`
	Detail   string `json:"detail" form:"detail"`
	Slug     string `json:"slug" form:"slug"`
	Price    int64  `json:"price" form:"price"`
	Stock    int64  `json:"stock" form:"stock"`
	Category string `json:"category" form:"category"`
	Image    string `json:"image" form:"image"`
}

type InserProductRequest struct {
	AdministratorID uint                  `json:"administrator_id" form:"administrator_id"`
	Slug            string                `json:"slug" form:"slug"`
	Name            string                `json:"name" form:"name"`
	Detail          string                `json:"detail" form:"detail"`
	Price           int64                 `json:"price" form:"price"`
	Stock           int64                 `json:"stock" form:"stock"`
	Category        string                `json:"category" form:"category"`
	Image           *multipart.FileHeader `json:"image" form:"image"`
}

type ImageProdukRequest struct {
	Image string `json:"image" form:"image" validate:"required"`
}

type InserProductResponse struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Detail   string `json:"detail"`
	Price    int64  `json:"price"`
	Stock    int64  `json:"stock"`
	Category string `json:"category"`
	Image    string `json:"image"`
}

type ProductResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Detail   string `json:"detail"`
	Price    int64  `json:"price"`
	Stock    int64  `json:"stock"`
	Category string `json:"category"`
	Image    string `json:"image"`
}
