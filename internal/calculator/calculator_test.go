package calculator

import (
	"errors"
	"testing"

	"finance/internal/domain"
)

func TestCalculateTotalBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		plan    domain.FinancialPlan
		want    domain.Money
		wantErr error
	}{
		{
			name: "sums account balances",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
					{Balance: rub(250)},
				},
			},
			want: rub(350),
		},
		{
			name: "empty plan returns zero RUB",
			want: rub(0),
		},
		{
			name: "rejects mixed currencies",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
					{Balance: domain.NewUSD(100_00)},
				},
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateTotalBalance(tt.plan)
			assertMoneyResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestCalculateMonthlyDebtPayments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		plan    domain.FinancialPlan
		want    domain.Money
		wantErr error
	}{
		{
			name: "sums monthly debt payments",
			plan: domain.FinancialPlan{
				Debts: []domain.Debt{
					{MonthlyPayment: rub(100)},
					{MonthlyPayment: rub(50)},
				},
			},
			want: rub(150),
		},
		{
			name: "empty plan returns zero RUB",
			want: rub(0),
		},
		{
			name: "rejects mixed currencies",
			plan: domain.FinancialPlan{
				Debts: []domain.Debt{
					{MonthlyPayment: rub(100)},
					{MonthlyPayment: domain.NewEUR(100_00)},
				},
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateMonthlyDebtPayments(tt.plan)
			assertMoneyResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestCalculateTotalIncome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		plan    domain.FinancialPlan
		want    domain.Money
		wantErr error
	}{
		{
			name: "sums recurring income",
			plan: domain.FinancialPlan{
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(200)},
					{Amount: rub(300)},
				},
			},
			want: rub(500),
		},
		{
			name: "empty plan returns zero RUB",
			want: rub(0),
		},
		{
			name: "rejects mixed currencies",
			plan: domain.FinancialPlan{
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(200)},
					{Amount: domain.NewUSD(300_00)},
				},
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateTotalIncome(tt.plan)
			assertMoneyResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestCalculateRequiredRecurringExpenses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		plan    domain.FinancialPlan
		want    domain.Money
		wantErr error
	}{
		{
			name: "sums only required recurring expenses",
			plan: domain.FinancialPlan{
				RecurringExpenses: []domain.RecurringExpense{
					{Amount: rub(100), Type: domain.ExpenseTypeRequired},
					{Amount: rub(250), Type: domain.ExpenseTypeOptional},
					{Amount: rub(50), Type: domain.ExpenseTypeRequired},
				},
			},
			want: rub(150),
		},
		{
			name: "empty plan returns zero RUB",
			want: rub(0),
		},
		{
			name: "rejects mixed currencies among required expenses",
			plan: domain.FinancialPlan{
				RecurringExpenses: []domain.RecurringExpense{
					{Amount: rub(100), Type: domain.ExpenseTypeRequired},
					{Amount: domain.NewUSD(250_00), Type: domain.ExpenseTypeOptional},
					{Amount: domain.NewUSD(50_00), Type: domain.ExpenseTypeRequired},
				},
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateRequiredRecurringExpenses(tt.plan)
			assertMoneyResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestCalculateFreeMoney(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		plan    domain.FinancialPlan
		want    domain.Money
		wantErr error
	}{
		{
			name: "calculates balance plus income minus required expenses and debt payments",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
				},
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: rub(300)},
				},
				RecurringExpenses: []domain.RecurringExpense{
					{Amount: rub(70), Type: domain.ExpenseTypeRequired},
					{Amount: rub(20), Type: domain.ExpenseTypeOptional},
				},
				Debts: []domain.Debt{
					{MonthlyPayment: rub(50)},
				},
			},
			want: rub(280),
		},
		{
			name: "empty plan returns zero RUB",
			want: rub(0),
		},
		{
			name: "uses actual plan currency when some groups are empty",
			plan: domain.FinancialPlan{
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: domain.NewUSD(300_00)},
				},
			},
			want: domain.NewUSD(300_00),
		},
		{
			name: "rejects mixed currencies across totals",
			plan: domain.FinancialPlan{
				Accounts: []domain.Account{
					{Balance: rub(100)},
				},
				RecurringIncomes: []domain.RecurringIncome{
					{Amount: domain.NewUSD(300_00)},
				},
			},
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateFreeMoney(tt.plan)
			assertMoneyResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func rub(majorUnits int64) domain.Money {
	return domain.NewMoneyFromMajorUnits(majorUnits, domain.CurrencyRUB)
}

func assertMoneyResult(t *testing.T, got domain.Money, err error, want domain.Money, wantErr error) {
	t.Helper()

	if wantErr != nil {
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected error %v, got %v", wantErr, err)
		}
		return
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
