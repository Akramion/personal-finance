package calculator

import "finance/internal/domain"

func CalculateTotalIncome(plan domain.FinancialPlan) (domain.Money, error) {
	return calculateTotalIncome(plan, domain.CurrencyRUB)
}

func calculateTotalIncome(plan domain.FinancialPlan, defaultCurrency domain.Currency) (domain.Money, error) {
	incomes := make([]domain.Money, 0, len(plan.RecurringIncomes))

	for _, income := range plan.RecurringIncomes {
		incomes = append(incomes, income.Amount)
	}

	return domain.SumMoneyOrZero(defaultCurrency, incomes...)
}
