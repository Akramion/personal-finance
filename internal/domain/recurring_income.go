package domain

type IncomeType string

const (
	IncomeTypeSalary IncomeType = "salary"
	IncomeTypeBonus  IncomeType = "bonus"
	IncomeTypeGift   IncomeType = "gift"
	IncomeTypeOther  IncomeType = "other"
)

type RecurringIncome struct {
	ID          string
	UserID      string
	AccountID   string
	Amount      Money
	Type        IncomeType
	Description string
	PaymentDay  int
}
