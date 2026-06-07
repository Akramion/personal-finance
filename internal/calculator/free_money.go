package calculator

import "finance/internal/domain"

func CalculateFreeMoney(plan domain.FinancialPlan) (domain.Money, error) {
	currency, err := planCurrency(plan)
	if err != nil {
		return domain.Money{}, err
	}

	totalBalance, err := calculateTotalBalance(plan, currency)
	if err != nil {
		return domain.Money{}, err
	}

	totalIncome, err := calculateTotalIncome(plan, currency)
	if err != nil {
		return domain.Money{}, err
	}

	totalExpenses, err := calculateRequiredRecurringExpenses(plan, currency)
	if err != nil {
		return domain.Money{}, err
	}

	debtPayments, err := calculateMonthlyDebtPayments(plan, currency)
	if err != nil {
		return domain.Money{}, err
	}

	freeMoney, err := totalBalance.Add(totalIncome)
	if err != nil {
		return domain.Money{}, err
	}

	freeMoney, err = freeMoney.Sub(totalExpenses)
	if err != nil {
		return domain.Money{}, err
	}

	freeMoney, err = freeMoney.Sub(debtPayments)
	if err != nil {
		return domain.Money{}, err
	}

	return freeMoney, nil
}

func planCurrency(plan domain.FinancialPlan) (domain.Currency, error) {
	var currency domain.Currency
	var hasCurrency bool

	for _, account := range plan.Accounts {
		if err := setCurrency(&currency, &hasCurrency, account.Balance); err != nil {
			return "", err
		}
	}

	for _, income := range plan.RecurringIncomes {
		if err := setCurrency(&currency, &hasCurrency, income.Amount); err != nil {
			return "", err
		}
	}

	for _, expense := range plan.RecurringExpenses {
		if expense.Type != domain.ExpenseTypeRequired {
			continue
		}

		if err := setCurrency(&currency, &hasCurrency, expense.Amount); err != nil {
			return "", err
		}
	}

	for _, debt := range plan.Debts {
		if err := setCurrency(&currency, &hasCurrency, debt.MonthlyPayment); err != nil {
			return "", err
		}
	}

	if !hasCurrency {
		return domain.CurrencyRUB, nil
	}

	return currency, nil
}

func setCurrency(currency *domain.Currency, hasCurrency *bool, money domain.Money) error {
	if !*hasCurrency {
		*currency = money.Currency()
		*hasCurrency = true
		return nil
	}

	if *currency != money.Currency() {
		return domain.ErrCurrencyMismatch
	}

	return nil
}
