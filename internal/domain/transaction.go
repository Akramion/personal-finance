package domain

import "time"

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID          string
	UserID      string
	AccountID   string
	Amount      Money
	Type        TransactionType
	Category    string
	Description string
	OccurredAt  time.Time
}
