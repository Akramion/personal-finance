package domain

type ExpenseType string

const (
	ExpenseTypeRequired ExpenseType = "required"
	ExpenseTypeOptional ExpenseType = "optional"
)

type ExpenseCategory string

const (
	ExpenseCategoryRent          ExpenseCategory = "rent"
	ExpenseCategoryFood          ExpenseCategory = "food"
	ExpenseCategoryTransport     ExpenseCategory = "transport"
	ExpenseCategoryEntertainment ExpenseCategory = "entertainment"
	ExpenseCategoryHealth        ExpenseCategory = "health"
	ExpenseCategoryOther         ExpenseCategory = "other"
)

type RecurringExpense struct {
	ID          string
	UserID      string
	AccountID   string
	Amount      Money
	Type        ExpenseType
	Category    ExpenseCategory
	Description string
	PaymentDay  int
}
