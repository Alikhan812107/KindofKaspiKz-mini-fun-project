package postgres

import (
	"database/sql"
	"transaction-service/internal/domain"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(p *domain.Payment) error {
	_, err := r.db.Exec(
		`INSERT INTO payments (id, order_id, transaction_id, amount, status) VALUES ($1,$2,$3,$4,$5)`,
		p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status,
	)

	return err
}

func (r *PaymentRepo) GetByOrderID(orderID string) (*domain.Payment, error) {
	row := r.db.QueryRow(
		`SELECT id, order_id, transaction_id, amount, status FROM payments WHERE order_id=$1 LIMIT 1`, orderID,
	)
	p := &domain.Payment{}
	err := row.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	if err != nil {
		return nil, err
	}
	return p, nil
}
