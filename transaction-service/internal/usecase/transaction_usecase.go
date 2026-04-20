package usecase

import (
	"errors"
	"transaction-service/internal/domain"
	"transaction-service/internal/repository"

	"github.com/google/uuid"
)

type PaymentUsecase struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(r repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repo: r}
}

func (uc *PaymentUsecase) CreatePayment(orderID string, amount int64) (*domain.Payment, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	status := "Authorized"
	if amount > 100000 {
		status = "Declined"
	}

	p := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
		Status:        status,
	}

	if err := uc.repo.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (uc *PaymentUsecase) GetByOrderID(orderID string) (*domain.Payment, error) {
	return uc.repo.GetByOrderID(orderID)
}
