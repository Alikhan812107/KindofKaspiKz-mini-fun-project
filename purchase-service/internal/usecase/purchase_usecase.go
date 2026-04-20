package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"purchase-service/internal/domain"
	"purchase-service/internal/repository"

	"github.com/google/uuid"
)

type PurchaseUsecase struct {
	repo   repository.PurchaseRepository
	client *http.Client
}

func NewPurchaseUsecase(repo repository.PurchaseRepository) *PurchaseUsecase {
	return &PurchaseUsecase{
		repo:   repo,
		client: &http.Client{Timeout: 5 * time.Second},
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

	body, _ := json.Marshal(map[string]interface{}{
		"order_id": p.ID,
		"amount":   p.Amount,
	})

	resp, err := uc.client.Post("http://localhost:8081/payments", "application/json", bytes.NewBuffer(body))
	if err != nil {
		uc.repo.UpdateStatus(p.ID, "Failed")
		p.Status = "Failed"
		return p, errors.New("payment service unavailable")
	}
	defer resp.Body.Close()

	var txResp struct {
		Status string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&txResp)

	if txResp.Status == "Authorized" {
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
