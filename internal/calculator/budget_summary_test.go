package calculator

import (
	"errors"
	"testing"

	"finance/internal/domain"
)

func TestCalculateBudgetSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		plan              domain.FinancialPlan
		goalContributions domain.Money
		want              BudgetSummary
		wantErr           error
	}{
		{
			name: "positive free money",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
				},
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(400)},
				},
				RecurringExpenses: []domain.RecurringExpense{
					{Amount: rub(150), Type: domain.ExpenseTypeRequired},
					{Amount: rub(80), Type: domain.ExpenseTypeOptional},
				},
				Debts: []domain.Debt{
					{MonthlyPayment: rub(50)},
				},
			},
			goalContributions: rub(100),
			want: BudgetSummary{
				CurrentBalance:        rub(100),
				ExpectedIncome:        rub(400),
				RequiredExpenses:      rub(150),
				DebtPayments:          rub(50),
				GoalContributions:     rub(100),
				FreeMoney:             rub(200),
				OptionalSpendingLimit: rub(200),
			},
		},
		{
			name: "zero free money",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
				},
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(200)},
				},
				RecurringExpenses: []domain.RecurringExpense{
					{Amount: rub(100), Type: domain.ExpenseTypeRequired},
				},
				Debts: []domain.Debt{
					{MonthlyPayment: rub(50)},
				},
			},
			goalContributions: rub(150),
			want: BudgetSummary{
				CurrentBalance:        rub(100),
				ExpectedIncome:        rub(200),
				RequiredExpenses:      rub(100),
				DebtPayments:          rub(50),
				GoalContributions:     rub(150),
				FreeMoney:             rub(0),
				OptionalSpendingLimit: rub(0),
			},
		},
		{
			name: "negative free money",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(50)},
				},
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(100)},
				},
				RecurringExpenses: []domain.RecurringExpense{
					{Amount: rub(120), Type: domain.ExpenseTypeRequired},
				},
				Debts: []domain.Debt{
					{MonthlyPayment: rub(40)},
				},
			},
			goalContributions: rub(20),
			want: BudgetSummary{
				CurrentBalance:        rub(50),
				ExpectedIncome:        rub(100),
				RequiredExpenses:      rub(120),
				DebtPayments:          rub(40),
				GoalContributions:     rub(20),
				FreeMoney:             rub(-30),
				OptionalSpendingLimit: rub(0),
			},
		},
		{
			name: "currency mismatch",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
				},
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(100)},
				},
			},
			goalContributions: domain.NewUSD(10_00),
			wantErr:           domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateBudgetSummary(tt.plan, tt.goalContributions)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}
