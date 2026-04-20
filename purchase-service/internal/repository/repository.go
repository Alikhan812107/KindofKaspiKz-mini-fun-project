package repository

import "purchase-service/internal/domain"

type PurchaseRepository interface {
	Create(p *domain.Purchase) error
	GetByID(id string) (*domain.Purchase, error)
	UpdateStatus(id, status string) error
	GetByCustomerID(customerID string) ([]*domain.Purchase, error)
}
