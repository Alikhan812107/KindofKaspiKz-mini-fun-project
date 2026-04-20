package postgres

import (
	"database/sql"
	"purchase-service/internal/domain"
)

type PurchaseRepo struct {
	db *sql.DB
}

func NewPurchaseRepo(db *sql.DB) *PurchaseRepo {
	return &PurchaseRepo{db: db}
}

func (r *PurchaseRepo) Create(p *domain.Purchase) error {
	_, err := r.db.Exec(
		`INSERT INTO purchases (id, customer_id, item_name, amount, status, created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		p.ID, p.CustomerID, p.ItemName, p.Amount, p.Status, p.CreatedAt,
	)
	return err
}

func (r *PurchaseRepo) GetByID(id string) (*domain.Purchase, error) {
	row := r.db.QueryRow(
		`SELECT id, customer_id, item_name, amount, status, created_at FROM purchases WHERE id=$1`, id,
	)
	p := &domain.Purchase{}
	err := row.Scan(&p.ID, &p.CustomerID, &p.ItemName, &p.Amount, &p.Status, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PurchaseRepo) UpdateStatus(id, status string) error {
	_, err := r.db.Exec(`UPDATE purchases SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *PurchaseRepo) GetByCustomerID(customerID string) ([]*domain.Purchase, error) {
	rows, err := r.db.Query(
		`SELECT id, customer_id, item_name, amount, status, created_at FROM purchases WHERE customer_id=$1`, customerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*domain.Purchase
	for rows.Next() {
		p := &domain.Purchase{}
		if err := rows.Scan(&p.ID, &p.CustomerID, &p.ItemName, &p.Amount, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}
