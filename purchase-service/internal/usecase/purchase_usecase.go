package usecase

import (
	"context"
	"errors"
	"time"

	"purchase-service/internal/domain"
	"purchase-service/internal/repository"
	pb "purchase-service/pkg/pb/payment"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PurchaseUsecase struct {
	repo          repository.PurchaseRepository
	paymentClient pb.PaymentServiceClient
}

func NewPurchaseUsecase(repo repository.PurchaseRepository, paymentClient pb.PaymentServiceClient) *PurchaseUsecase {
	return &PurchaseUsecase{
		repo:          repo,
		paymentClient: paymentClient,
	}
}

func (uc *PurchaseUsecase) CreatePurchase(customerID, item string, amount int64) (*domain.Purchase, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	p := &domain.Purchase{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   item,
		Amount:     amount,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}

	if err := uc.repo.Create(p); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := uc.paymentClient.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: p.ID,
		Amount:  p.Amount,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Unavailable {
			uc.repo.UpdateStatus(p.ID, "Failed")
			p.Status = "Failed"
			return p, errors.New("payment service unavailable")
		}
		uc.repo.UpdateStatus(p.ID, "Failed")
		p.Status = "Failed"
		return p, err
	}

	if resp.Status == "Authorized" {
		p.Status = "Paid"
	} else {
		p.Status = "Failed"
	}

	uc.repo.UpdateStatus(p.ID, p.Status)
	return p, nil
}

func (uc *PurchaseUsecase) GetByID(id string) (*domain.Purchase, error) {
	return uc.repo.GetByID(id)
}

func (uc *PurchaseUsecase) GetByCustomerID(customerID string) ([]*domain.Purchase, error) {
	return uc.repo.GetByCustomerID(customerID)
}

func (uc *PurchaseUsecase) CancelPurchase(id string) (*domain.Purchase, error) {
	p, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if p.Status == "Paid" {
		return nil, errors.New("paid orders cannot be cancelled")
	}
	if p.Status == "Cancelled" {
		return nil, errors.New("order already cancelled")
	}
	if err := uc.repo.UpdateStatus(id, "Cancelled"); err != nil {
		return nil, err
	}
	p.Status = "Cancelled"
	return p, nil
}
