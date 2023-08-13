package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type CartDetailRepository struct {
	db *gorm.DB
}

func NewCartDetailRepository(db *gorm.DB) models.CartDetailRepository {
	return &CartDetailRepository{db}
}

func (r *CartDetailRepository) InsertOne(req *models.CartDetail) (*models.CartDetail, error) {
	err := r.db.Create(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *CartDetailRepository) FindOne(userID, productID uint) (res *models.CartDetail, err error) {
	var CartDetail models.CartDetail

	err = r.db.Model(&CartDetail).Where("user_id = ? AND product_id = ?", userID, productID).First(&CartDetail).Error
	if err != nil {
		return nil, err
	}

	return &CartDetail, nil
}

func (r *CartDetailRepository) FindAll(userID uint) (*[]models.CartDetail, error) {
	var res []models.CartDetail
	err := r.db.Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *CartDetailRepository) UpdateOne(req *models.CartDetail, id uint) (res *models.CartDetail, err error) {
	err = r.db.Where("id = ?", id).Updates(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *CartDetailRepository) DeleteOne(req *models.CartDetail) error {
	err := r.db.Delete(&req).Error
	if err != nil {
		return err
	}

	return nil
}
