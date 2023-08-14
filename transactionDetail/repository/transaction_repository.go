package repository

import (
	"context"
	"fiberStore/models"

	"gorm.io/gorm"
)

type TransactionDetailRepository struct {
	db *gorm.DB
}

func NewTransactionDetailRepository(db *gorm.DB) models.TransactionDetailRepository {
	return &TransactionDetailRepository{db}
}

func (tr *TransactionDetailRepository) BeginTx(ctx context.Context) *gorm.DB {
	return tr.db.Begin()
}

func (r *TransactionDetailRepository) InsertOne(req *models.TransactionDetail) (*models.TransactionDetail, error) {
	err := r.db.Create(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *TransactionDetailRepository) FindOne(id uint) (res *models.TransactionDetail, err error) {
	var TransactionDetail models.TransactionDetail

	err = r.db.Model(&TransactionDetail).Where("id = ?", id).First(&TransactionDetail).Error
	if err != nil {
		return nil, err
	}

	return &TransactionDetail, nil
}

func (r *TransactionDetailRepository) FindAll(userID uint) (*[]models.TransactionDetail, error) {
	var res []models.TransactionDetail
	err := r.db.Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *TransactionDetailRepository) UpdateOne(req *models.TransactionDetail, id uint) (res *models.TransactionDetail, err error) {
	err = r.db.Where("id = ?", id).Updates(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *TransactionDetailRepository) DeleteOne(req *models.TransactionDetail) error {
	err := r.db.Delete(&req).Error
	if err != nil {
		return err
	}

	return nil
}
