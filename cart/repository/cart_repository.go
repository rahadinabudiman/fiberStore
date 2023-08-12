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

func (r *cartRepository) FindOne(userID, productID uint) (res *models.Cart, err error) {
	var cart models.Cart

	err = r.db.Model(&cart).Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepository) FindAll(userID uint) (*[]models.Cart, error) {
	var res []models.Cart
	err := r.db.Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *cartRepository) UpdateOne(req *models.Cart, id uint) (res *models.Cart, err error) {
	err = r.db.Where("id = ?", id).Updates(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *cartRepository) DeleteOne(req *models.Cart) error {
	err := r.db.Delete(&req).Error
	if err != nil {
		return err
	}

	return nil
}
