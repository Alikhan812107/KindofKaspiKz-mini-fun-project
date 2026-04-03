package domain

type Transaction struct {
	ID            string
	PurchaseID    string
	TransactionID string
	Amount        int64
	Status        string // "Authorized" or "Declined"
}
