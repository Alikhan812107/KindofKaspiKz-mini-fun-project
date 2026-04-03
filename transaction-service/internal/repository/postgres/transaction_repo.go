package postgres

import (
	"database/sql"
	"transaction-service/internal/domain"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, purchase_id, transaction_id, amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		tx.ID,
		tx.PurchaseID,
		tx.TransactionID,
		tx.Amount,
		tx.Status,
	)

	return err
}

func (r *TransactionRepo) GetByPurchaseID(purchaseID string) (*domain.Transaction, error) {
	query := `
		SELECT id, purchase_id, transaction_id, amount, status
		FROM transactions
		WHERE purchase_id = $1
		LIMIT 1
	`

	row := r.db.QueryRow(query, purchaseID)

	tx := &domain.Transaction{}

	err := row.Scan(
		&tx.ID,
		&tx.PurchaseID,
		&tx.TransactionID,
		&tx.Amount,
		&tx.Status,
	)

	if err != nil {
		return nil, err
	}

	return tx, nil
}
