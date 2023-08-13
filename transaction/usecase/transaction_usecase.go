package usecase

import (
	"context"
	"errors"
	"fiberStore/models"
	"time"
)

type transactionUsecase struct {
	TransactionRepository models.TransactionRepository
	CartRepository        models.CartRepository
	CartDetailRepository  models.CartDetailRepository
	ProductRepository     models.ProductRepository
	UserRepository        models.UserRepository
	UserAmountRepository  models.UserAmountRepository
	contextTimeout        time.Duration
}

func NewTransactionUsecase(TransactionRepository models.TransactionRepository, CartRepository models.CartRepository, CartDetailRepository models.CartDetailRepository, ProductRepository models.ProductRepository, UserRepository models.UserRepository, UserAmountRepository models.UserAmountRepository, contextTimeout time.Duration) models.TransactionUsecase {
	return &transactionUsecase{
		TransactionRepository: TransactionRepository,
		CartRepository:        CartRepository,
		CartDetailRepository:  CartDetailRepository,
		ProductRepository:     ProductRepository,
		UserRepository:        UserRepository,
		UserAmountRepository:  UserAmountRepository,
		contextTimeout:        contextTimeout,
	}
}

func (tu *transactionUsecase) InsertOne(ctx context.Context, req *models.Transaction) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	// Memulai transaksi
	tx := tu.TransactionRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mencari data cart dan cart details
	cart, err := tu.CartRepository.FindOne(req.UserID)
	if err != nil {
		return nil, err
	}

	cartDetails, err := tu.CartDetailRepository.FindAll(cart.UserID)
	if err != nil {
		return nil, err
	}

	// Menghitung total harga dari cart details
	var totalPrice int
	for _, detail := range *cartDetails {
		product, err := tu.ProductRepository.FindOne(int(detail.ProductID))
		if err != nil {
			return nil, err
		}
		totalPrice += int(product.Price) * detail.Quantity
	}

	// Memastikan saldo pengguna cukup
	userAmount, err := tu.UserAmountRepository.FindOne(cart.UserID)
	if err != nil {
		return nil, err
	}

	if userAmount.Amount < float64(totalPrice) {
		return nil, errors.New("insufficient amount")
	}

	// Membuat transaksi baru
	newTransaction := &models.Transaction{
		UserID:     req.UserID,
		CartID:     cart.ID,
		TotalPrice: totalPrice,
	}

	// Menyimpan transaksi baru
	insertedTransaction, err := tu.TransactionRepository.InsertOne(newTransaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Mengurangi stok produk dan menghapus cart details
	for _, detail := range *cartDetails {
		product, err := tu.ProductRepository.FindOne(int(detail.ProductID))
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if product.Stock < int64(detail.Quantity) {
			tx.Rollback()
			return nil, errors.New("product stock not enough")
		}

		product.Stock -= int64(detail.Quantity)
		_, err = tu.ProductRepository.UpdateOne(product)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		err = tu.CartDetailRepository.DeleteOne(&detail)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Mengurangi saldo pengguna
	userAmount.Amount -= float64(totalPrice)
	_, err = tu.UserAmountRepository.UpdateOne(userAmount, cart.UserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Menghapus cart
	err = tu.CartRepository.DeleteOne(cart)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaksi
	tx.Commit()

	return insertedTransaction, nil
}

func (tu *transactionUsecase) FindOne(ctx context.Context, id uint) (*models.Transaction, error) {
	return nil, nil
}

func (tu *transactionUsecase) FindAll(ctx context.Context, userID uint) (*[]models.Transaction, error) {
	return nil, nil
}
