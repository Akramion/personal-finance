package calculator

import "finance/internal/domain"

type BudgetSummary struct {
	CurrentBalance        domain.Money
	ExpectedIncome        domain.Money
	RequiredExpenses      domain.Money
	DebtPayments          domain.Money
	GoalContributions     domain.Money
	FreeMoney             domain.Money
	OptionalSpendingLimit domain.Money
}

func CalculateBudgetSummary(plan domain.FinancialPlan, goalContributions domain.Money) (BudgetSummary, error) {
	currency, err := budgetSummaryCurrency(plan, goalContributions)
	if err != nil {
		return BudgetSummary{}, err
	}

	currentBalance, err := calculateTotalBalance(plan, currency)
	if err != nil {
		return BudgetSummary{}, err
	}

	expectedIncome, err := calculateTotalIncome(plan, currency)
	if err != nil {
		return BudgetSummary{}, err
	}

	requiredExpenses, err := calculateRequiredRecurringExpenses(plan, currency)
	if err != nil {
		return BudgetSummary{}, err
	}

	debtPayments, err := calculateMonthlyDebtPayments(plan, currency)
	if err != nil {
		return BudgetSummary{}, err
	}

	freeMoney, err := currentBalance.Add(expectedIncome)
	if err != nil {
		return BudgetSummary{}, err
	}

	freeMoney, err = freeMoney.Sub(requiredExpenses)
	if err != nil {
		return BudgetSummary{}, err
	}

	freeMoney, err = freeMoney.Sub(debtPayments)
	if err != nil {
		return BudgetSummary{}, err
	}

	freeMoney, err = freeMoney.Sub(goalContributions)
	if err != nil {
		return BudgetSummary{}, err
	}

	optionalSpendingLimit := freeMoney
	if freeMoney.MinorUnits() < 0 {
		optionalSpendingLimit = domain.ZeroMoney(currency)
	}

	return BudgetSummary{
		CurrentBalance:        currentBalance,
		ExpectedIncome:        expectedIncome,
		RequiredExpenses:      requiredExpenses,
		DebtPayments:          debtPayments,
		GoalContributions:     goalContributions,
		FreeMoney:             freeMoney,
		OptionalSpendingLimit: optionalSpendingLimit,
	}, nil
}

func budgetSummaryCurrency(plan domain.FinancialPlan, goalContributions domain.Money) (domain.Currency, error) {
	currency, err := planCurrency(plan)
	if err != nil {
		return "", err
	}

	if currency != goalContributions.Currency() {
		if hasPlanMoney(plan) {
			return "", domain.ErrCurrencyMismatch
		}

		return goalContributions.Currency(), nil
	}

	return currency, nil
}

func hasPlanMoney(plan domain.FinancialPlan) bool {
	return len(plan.Accounts) > 0 ||
		len(plan.RecurringIncomes) > 0 ||
		hasRequiredExpenses(plan) ||
		len(plan.Debts) > 0
}

func hasRequiredExpenses(plan domain.FinancialPlan) bool {
	for _, expense := range plan.RecurringExpenses {
		if expense.Type == domain.ExpenseTypeRequired {
			return true
		}
	}

	return false
}
