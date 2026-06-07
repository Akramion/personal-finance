package calculator

import "finance/internal/domain"

func CalculateTotalBalance(plan domain.FinancialPlan) (domain.Money, error) {
	return calculateTotalBalance(plan, domain.CurrencyRUB)
}

func calculateTotalBalance(plan domain.FinancialPlan, defaultCurrency domain.Currency) (domain.Money, error) {
	balances := make([]domain.Money, 0, len(plan.Accounts))

	for _, account := range plan.Accounts {
		balances = append(balances, account.Balance)
	}

	return domain.SumMoneyOrZero(defaultCurrency, balances...)
}
