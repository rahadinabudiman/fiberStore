package repository

import (
	"context"
	"fiberStore/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	InsertOne(ctx context.Context, req *models.User) (*models.User, error)
	FindOneById(ctx context.Context, id int) (*models.User, error)
	FindAll(ctx context.Context, page, limit int, search string) (*[]models.User, int, error)
	UpdateOne(ctx context.Context, req *models.User) (*models.User, error)
	DeleteOne(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) InsertOne(ctx context.Context, req *models.User) (*models.User, error) {
	return nil, nil
}

func (ur *userRepository) FindOneById(ctx context.Context, id int) (*models.User, error) {
	return nil, nil
}

func (ur *userRepository) FindAll(ctx context.Context, page, limit int, search string) (*[]models.User, int, error) {
	return nil, 0, nil
}

func (ur *userRepository) UpdateOne(ctx context.Context, req *models.User) (*models.User, error) {
	return nil, nil
}

func (ur *userRepository) DeleteOne(ctx context.Context, id int) error {
	return nil
}
