package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) models.CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) InsertOne(req *models.Cart) (*models.Cart, error) {
	err := r.db.Create(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *cartRepository) FindOne(userID uint) (res *models.Cart, err error) {
	var cart models.Cart

	err = r.db.Model(&cart).Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) DeleteOne(req *models.Cart) error {
	err := r.db.Delete(&req).Error
	if err != nil {
		return err
	}

	return nil
}
