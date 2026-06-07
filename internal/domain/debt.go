package domain

import "time"

type DebtType string

const (
	DebtTypeCreditCard  DebtType = "credit_card"
	DebtTypeLoan        DebtType = "loan"
	DebtTypeInstallment DebtType = "installment"
	DebtTypeOther       DebtType = "other"
)

type Debt struct {
	ID              string
	UserID          string
	Name            string
	Type            DebtType
	TotalAmount     Money
	RemainingAmount Money
	MonthlyPayment  Money
	PaymentDay      int
	CreatedAt       time.Time
}
