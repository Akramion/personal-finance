package calculator

import "finance/internal/domain"

func CalculateRequiredRecurringExpenses(plan domain.FinancialPlan) (domain.Money, error) {
	return calculateRequiredRecurringExpenses(plan, domain.CurrencyRUB)
}

func calculateRequiredRecurringExpenses(plan domain.FinancialPlan, defaultCurrency domain.Currency) (domain.Money, error) {
	requiredExpenses := make([]domain.Money, 0, len(plan.RecurringExpenses))

	for _, expense := range plan.RecurringExpenses {
		if expense.Type != domain.ExpenseTypeRequired {
			continue
		}

		requiredExpenses = append(requiredExpenses, expense.Amount)
	}

	return domain.SumMoneyOrZero(defaultCurrency, requiredExpenses...)
}
