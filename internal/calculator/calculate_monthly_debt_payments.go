package calculator

import "finance/internal/domain"

func CalculateMonthlyDebtPayments(plan domain.FinancialPlan) (domain.Money, error) {
	return calculateMonthlyDebtPayments(plan, domain.CurrencyRUB)
}

func calculateMonthlyDebtPayments(plan domain.FinancialPlan, defaultCurrency domain.Currency) (domain.Money, error) {
	monthlyPayments := make([]domain.Money, 0, len(plan.Debts))

	for _, debt := range plan.Debts {
		monthlyPayments = append(monthlyPayments, debt.MonthlyPayment)
	}

	return domain.SumMoneyOrZero(defaultCurrency, monthlyPayments...)
}
