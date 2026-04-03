package usecase

import (
	"errors"
	"transaction-service/internal/domain"
	"transaction-service/internal/repository"

	"github.com/google/uuid"
)

type TransactionUsecase struct {
	repo repository.TransactionRepository
}

func NewTransactionUsecase(r repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{repo: r}
}

func (uc *TransactionUsecase) CreateTransaction(purchaseID string, amount int64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	status := "Authorized"

	if amount > 100000 {
		status = "Declined"
	}

	tx := &domain.Transaction{
		ID:            uuid.New().String(),
		PurchaseID:    purchaseID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
	}

	err := uc.repo.Create(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (uc *TransactionUsecase) GetByPurchaseID(purchaseID string) (*domain.Transaction, error) {
	return uc.repo.GetByPurchaseID(purchaseID)
}
