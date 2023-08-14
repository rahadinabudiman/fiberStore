package usecase

import (
	"context"
	"errors"
	"fiberStore/dtos"
	"fiberStore/models"
	"time"
)

type TransactionDetailUsecase struct {
	TransactionDetailRepository models.TransactionDetailRepository
	TransactionRepository       models.TransactionRepository
	CartRepository              models.CartRepository
	CartDetailRepository        models.CartDetailRepository
	ProductRepository           models.ProductRepository
	UserRepository              models.UserRepository
	UserAmountRepository        models.UserAmountRepository
	contextTimeout              time.Duration
}

func NewTransactionDetailUsecase(TransactionDetailRepository models.TransactionDetailRepository, TransactionRepository models.TransactionRepository, CartRepository models.CartRepository, CartDetailRepository models.CartDetailRepository, ProductRepository models.ProductRepository, UserRepository models.UserRepository, UserAmountRepository models.UserAmountRepository, contextTimeout time.Duration) models.TransactionDetailUsecase {
	return &TransactionDetailUsecase{
		TransactionDetailRepository: TransactionDetailRepository,
		TransactionRepository:       TransactionRepository,
		CartRepository:              CartRepository,
		CartDetailRepository:        CartDetailRepository,
		ProductRepository:           ProductRepository,
		UserRepository:              UserRepository,
		UserAmountRepository:        UserAmountRepository,
		contextTimeout:              contextTimeout,
	}
}

func (tu *TransactionDetailUsecase) InsertOne(ctx context.Context, req *models.TransactionDetail) (*dtos.InsertTransactionDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	// Memulai transaksi
	tx := tu.TransactionDetailRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mencari data user
	user, err := tu.UserRepository.FindOneById(int(req.UserID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Mencari data cart dan cart details
	cart, err := tu.CartRepository.FindOne(req.UserID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	cartDetails, err := tu.CartDetailRepository.FindAll(cart.UserID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Check Cart Exist or Create New One
	transaction, err := tu.TransactionRepository.FindOne(req.UserID)
	if err != nil {
		transaction = &models.Transaction{UserID: req.UserID}
		newTransaction, err := tu.TransactionRepository.InsertOne(transaction)
		if err != nil {
			return nil, errors.New("failed to add new transaction")
		}
		transaction = newTransaction
	}

	// Menghitung total harga dari cart details
	var totalPrice int64
	for _, detail := range *cartDetails {
		product, err := tu.ProductRepository.FindOne(int(detail.ProductID))
		if err != nil {
			return nil, errors.New("product not found")
		}
		totalPrice += product.Price * detail.Quantity
	}

	// Memastikan saldo pengguna cukup
	userAmount, err := tu.UserAmountRepository.FindOne(cart.UserID)
	if err != nil {
		return nil, errors.New("balance account not found")
	}

	if userAmount.Amount < float64(totalPrice) {
		return nil, errors.New("insufficient balance")
	}

	// Menghitung total harga dari cart details dan membuat transaksi detail
	var ProductResponse []dtos.ProductTransactionResponse
	for _, detail := range *cartDetails {
		product, err := tu.ProductRepository.FindOne(int(detail.ProductID))
		if err != nil {
			return nil, errors.New("product not found")
		}
		totalPricePerProduct := product.Price * detail.Quantity

		if product.Stock < int64(detail.Quantity) {
			tx.Rollback()
			return nil, errors.New("product stock not enough")
		}

		newTransactionDetail := &models.TransactionDetail{
			TransactionID: transaction.ID,
			UserID:        req.UserID,
			CartID:        cart.ID,
			ProductID:     detail.ProductID,
			Price:         product.Price,
			Quantity:      detail.Quantity,
			TotalPrice:    int64(totalPricePerProduct),
		}

		_, err = tu.TransactionDetailRepository.InsertOne(newTransactionDetail)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("failed to create TransactionDetail")
		}

		ProductResponse = append(ProductResponse, dtos.ProductTransactionResponse{
			Name:       product.Name,
			Price:      product.Price,
			Quantity:   detail.Quantity,
			TotalPrice: int64(totalPricePerProduct),
		})

		// Mengurangi stock produk
		product.Stock -= int64(detail.Quantity)
		_, err = tu.ProductRepository.UpdateOne(product)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("failed to update product")
		}

		// Menghapus Cart Detail
		err = tu.CartDetailRepository.DeleteOne(&detail)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("failed to delete cart detail")
		}
	}

	// Mengurangi saldo pengguna
	userAmount.Amount -= float64(totalPrice)
	_, err = tu.UserAmountRepository.UpdateOne(userAmount, cart.UserID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to update user balance account")
	}

	// Menghapus cart
	err = tu.CartRepository.DeleteOne(cart)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to delete cart")
	}

	// Menghapus ID Transaksi
	err = tu.TransactionRepository.DeleteOne(transaction)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("failed to delete transaction")
	}

	// Commit transaksi
	tx.Commit()

	// Membuat response
	res := &dtos.InsertTransactionDetailResponse{
		TransactionID: transaction.ID,
		Name:          user.Name,
		Produk:        ProductResponse,
		TotalPrice:    totalPrice,
	}

	return res, nil
}

func (tu *TransactionDetailUsecase) FindOne(ctx context.Context, id uint) (*models.TransactionDetail, error) {
	return nil, nil
}

func (tu *TransactionDetailUsecase) FindAll(ctx context.Context, userID uint) (*[]dtos.DetailTransactionDetailResponse, error) {
	_, cancel := context.WithTimeout(ctx, tu.contextTimeout)
	defer cancel()

	// Memulai transaksi
	tx := tu.TransactionDetailRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mencari data user
	user, err := tu.UserRepository.FindOneById(int(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	transactionDetail, err := tu.TransactionDetailRepository.FindAll(user.ID)
	if err != nil {
		return nil, errors.New("transaction detail not found")
	}

	var detailTransactionDetailResponses []dtos.DetailTransactionDetailResponse
	for _, detail := range *transactionDetail {
		product, err := tu.ProductRepository.FindOne(int(detail.ProductID))
		if err != nil {
			return nil, errors.New("product not found")
		}

		detailTransactionDetailResponse := dtos.DetailTransactionDetailResponse{
			ID:          detail.ID,
			CreatedAt:   detail.CreatedAt,
			ProductName: product.Name,
			Price:       detail.Price,
			Quantity:    detail.Quantity,
			TotalPrice:  detail.TotalPrice,
		}

		detailTransactionDetailResponses = append(detailTransactionDetailResponses, detailTransactionDetailResponse)
	}

	return &detailTransactionDetailResponses, nil
}
