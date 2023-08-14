package repository

import (
	"fiberStore/models"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) models.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (tr *transactionRepository) InsertOne(req *models.Transaction) (*models.Transaction, error) {
	err := tr.db.Create(req).Error
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (tr *transactionRepository) FindOne(userID uint) (res *models.Transaction, err error) {
	err = tr.db.Where("user_id = ?", userID).First(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (tr *transactionRepository) DeleteOne(req *models.Transaction) error {
	err := tr.db.Delete(req).Error
	if err != nil {
		return err
	}

	return nil
}
