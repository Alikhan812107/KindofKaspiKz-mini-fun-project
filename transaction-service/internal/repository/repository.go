package repository

import "transaction-service/internal/domain"

type TransactionRepository interface {
	Create(tx *domain.Transaction) error
	GetByPurchaseID(purchaseID string) (*domain.Transaction, error)
}
