package usecase

import (
	"fiberStore/models"
	"time"
)

type productUsecase struct {
	ProductRepository models.ProductRepository
	UserRepository    models.UserRepository
	contextTimeout    time.Duration
}
