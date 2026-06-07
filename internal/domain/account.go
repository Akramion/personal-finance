package domain

type AccountType string

const (
	AccountTypeCash AccountType = "cash"
	AccountTypeCard AccountType = "card"
	AccountTypeBank AccountType = "bank"
)

type Account struct {
	ID      string
	UserID  string
	Name    string
	Type    AccountType
	Balance Money
}
