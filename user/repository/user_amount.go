package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type userAmountRepository struct {
	db *gorm.DB
}

func NewUserAmountRepository(db *gorm.DB) models.UserAmountRepository {
	return &userAmountRepository{db}
}

func (uar *userAmountRepository) InsertOne(req *models.UserAmount) (*models.UserAmount, error) {
	err := uar.db.Create(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (uar *userAmountRepository) FindOne(id uint) (res *models.UserAmount, err error) {
	var userAmount *models.UserAmount

	err = uar.db.Model(&userAmount).Where("user_id = ?", id).First(&userAmount).Error
	if err != nil {
		return nil, err
	}

	return userAmount, nil
}

func (uar *userAmountRepository) UpdateOne(req *models.UserAmount, id uint) (res *models.UserAmount, err error) {
	var userAmount *models.UserAmount

	err = uar.db.Model(&userAmount).Where("user_id = ?", id).Updates(req).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (uar *userAmountRepository) DeleteOne(req *models.UserAmount) error {
	err := uar.db.Delete(req).Error
	if err != nil {
		return err
	}

	return nil
}
