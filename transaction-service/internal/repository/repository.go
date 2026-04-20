package repository

import "transaction-service/internal/domain"

type PaymentRepository interface {
	Create(p *domain.Payment) error
	GetByOrderID(orderID string) (*domain.Payment, error)
}
