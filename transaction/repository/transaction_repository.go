package repository

import (
	"context"
	"fiberStore/models"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) models.TransactionRepository {
	return &transactionRepository{db}
}

func (tr *transactionRepository) BeginTx(ctx context.Context) *gorm.DB {
	return tr.db.Begin()
}

func (r *transactionRepository) InsertOne(req *models.Transaction) (*models.Transaction, error) {
	err := r.db.Create(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *transactionRepository) FindOne(id uint) (res *models.Transaction, err error) {
	var transaction models.Transaction

	err = r.db.Model(&transaction).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) FindAll(userID uint) (*[]models.Transaction, error) {
	var res []models.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *transactionRepository) UpdateOne(req *models.Transaction, id uint) (res *models.Transaction, err error) {
	err = r.db.Where("id = ?", id).Updates(&req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *transactionRepository) DeleteOne(req *models.Transaction) error {
	err := r.db.Delete(&req).Error
	if err != nil {
		return err
	}

	return nil
}
