package models

import (
	"context"
	"fiberStore/dtos"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	AdministratorID uint   `gorm:"foreignKey:AdministratorRefer" json:"administrator_id" form:"administrator_id"`
	Administrator   User   `gorm:"foreignKey:AdministratorID" json:"-"`
	Slug            string `json:"slug" form:"slug"`
	Name            string `json:"name" form:"name"`
	Detail          string `json:"detail" form:"detail"`
	Price           int64  `json:"price" form:"price"`
	Stock           int64  `json:"stock" form:"stock"`
	Category        string `json:"category" form:"category"`
	Image           string `json:"image" form:"image"`
}

type ProductRepository interface {
	InsertOne(req *Product) (*Product, error)
	FindOne(id int) (*Product, error)
	FindOneBySlug(slug string) (*Product, error)
	FindByCategory(page, limit int, search string) (*[]Product, int, error)
	FindAll(page, limit int) (*[]Product, int, error)
	FindQueryAll(page, limit int, search string) (*[]Product, int, error)
	UpdateOne(req *Product) (*Product, error)
	DeleteOne(product *Product) error
}

type ProductUsecase interface {
	InsertOne(ctx context.Context, req *dtos.InserProductRequest, url string) (*dtos.InserProductResponse, error)
	FindAll(ctx context.Context, page, limit int) (*[]dtos.ProductResponse, int, error)
	FindAllByCategory(ctx context.Context, page, limit int, category string) (*[]dtos.ProductResponse, int, error)
	FindQueryAll(ctx context.Context, page, limit int, search string) (*[]dtos.ProductResponse, int, error)
	FindOne(ctx context.Context, id uint) (*dtos.ProductResponse, error)
	UpdateOne(ctx context.Context, req *dtos.UpdateProductRequest, id, AdministratorID uint) (*dtos.ProductResponse, error)
	DeleteOne(ctx context.Context, id, AdministratorID uint) error
}
