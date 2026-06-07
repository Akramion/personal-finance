package calculator

import (
	"errors"
	"testing"

	"finance/internal/domain"
)

func TestCalculateGoalPlanOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		goal          SavingsGoal
		budgetSummary BudgetSummary
		want          []GoalPlanOption
		wantErr       error
	}{
		{
			name: "returns comfortable balanced aggressive options",
			goal: SavingsGoal{
				Name:          "Emergency fund",
				TargetAmount:  rub(1000),
				CurrentAmount: rub(0),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			want: []GoalPlanOption{
				{
					Type:                GoalPlanOptionTypeComfortable,
					TargetAmount:        rub(1000),
					CurrentAmount:       rub(0),
					RemainingAmount:     rub(1000),
					MonthlyContribution: rub(50),
					Months:              20,
					IsAffordable:        true,
				},
				{
					Type:                GoalPlanOptionTypeBalanced,
					TargetAmount:        rub(1000),
					CurrentAmount:       rub(0),
					RemainingAmount:     rub(1000),
					MonthlyContribution: rub(80),
					Months:              13,
					IsAffordable:        true,
				},
				{
					Type:                GoalPlanOptionTypeAggressive,
					TargetAmount:        rub(1000),
					CurrentAmount:       rub(0),
					RemainingAmount:     rub(1000),
					MonthlyContribution: rub(100),
					Months:              10,
					IsAffordable:        true,
				},
			},
		},
		{
			name: "uses remaining amount when goal is partially completed",
			goal: SavingsGoal{
				Name:          "Laptop",
				TargetAmount:  rub(1000),
				CurrentAmount: rub(250),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			want: []GoalPlanOption{
				{
					Type:                GoalPlanOptionTypeComfortable,
					TargetAmount:        rub(1000),
					CurrentAmount:       rub(250),
					RemainingAmount:     rub(750),
					MonthlyContribution: rub(50),
					Months:              15,
					IsAffordable:        true,
				},
				{
					Type:                GoalPlanOptionTypeBalanced,
					TargetAmount:        rub(1000),
					CurrentAmount:       rub(250),
					RemainingAmount:     rub(750),
					MonthlyContribution: rub(80),
					Months:              10,
					IsAffordable:        true,
				},
				{
					Type:                GoalPlanOptionTypeAggressive,
					TargetAmount:        rub(1000),
					CurrentAmount:       rub(250),
					RemainingAmount:     rub(750),
					MonthlyContribution: rub(100),
					Months:              8,
					IsAffordable:        true,
				},
			},
		},
		{
			name: "returns reached result when current amount equals target amount",
			goal: SavingsGoal{
				Name:          "Trip",
				TargetAmount:  rub(500),
				CurrentAmount: rub(500),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(0)),
			want:          reachedOptions(rub(500), rub(500)),
		},
		{
			name: "returns reached result when current amount is greater than target amount",
			goal: SavingsGoal{
				Name:          "Trip",
				TargetAmount:  rub(500),
				CurrentAmount: rub(700),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			want:          reachedOptions(rub(500), rub(700)),
		},
		{
			name: "returns empty options when available monthly money is zero",
			goal: SavingsGoal{
				Name:          "Phone",
				TargetAmount:  rub(300),
				CurrentAmount: rub(0),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(0)),
			want:          []GoalPlanOption{},
		},
		{
			name: "returns error for goal currency mismatch",
			goal: SavingsGoal{
				Name:          "Phone",
				TargetAmount:  rub(300),
				CurrentAmount: domain.NewUSD(0),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr:       domain.ErrCurrencyMismatch,
		},
		{
			name: "returns error for budget currency mismatch",
			goal: SavingsGoal{
				Name:          "Phone",
				TargetAmount:  rub(300),
				CurrentAmount: rub(0),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(domain.NewUSD(100_00)),
			wantErr:       domain.ErrCurrencyMismatch,
		},
		{
			name: "returns error for zero target amount",
			goal: SavingsGoal{
				Name:          "Invalid",
				TargetAmount:  rub(0),
				CurrentAmount: rub(0),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr:       ErrInvalidGoalTarget,
		},
		{
			name: "returns error for negative target amount",
			goal: SavingsGoal{
				Name:          "Invalid",
				TargetAmount:  rub(-100),
				CurrentAmount: rub(0),
			},
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr:       ErrInvalidGoalTarget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateGoalPlanOptions(tt.goal, tt.budgetSummary)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d options, got %d: %+v", len(tt.want), len(got), got)
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("option %d: expected %+v, got %+v", i, tt.want[i], got[i])
				}
			}
		})
	}
}

func TestCalculateGoalDeadline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		goal          SavingsGoal
		months        int
		budgetSummary BudgetSummary
		want          GoalDeadlineResult
		wantErr       error
	}{
		{
			name: "calculates required monthly contribution",
			goal: SavingsGoal{
				Name:          "Emergency fund",
				TargetAmount:  rub(1200),
				CurrentAmount: rub(0),
			},
			months:        12,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			want: GoalDeadlineResult{
				TargetAmount:                 rub(1200),
				CurrentAmount:                rub(0),
				RemainingAmount:              rub(1200),
				Months:                       12,
				RequiredMonthlyContribution:  rub(100),
				AvailableMonthlyContribution: rub(100),
				IsAffordable:                 true,
				MonthlyShortage:              rub(0),
			},
		},
		{
			name: "rounds required monthly contribution up",
			goal: SavingsGoal{
				Name:          "Laptop",
				TargetAmount:  rub(1000),
				CurrentAmount: rub(0),
			},
			months:        3,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(400)),
			want: GoalDeadlineResult{
				TargetAmount:                 rub(1000),
				CurrentAmount:                rub(0),
				RemainingAmount:              rub(1000),
				Months:                       3,
				RequiredMonthlyContribution:  domain.NewMoney(333_34, domain.CurrencyRUB),
				AvailableMonthlyContribution: rub(400),
				IsAffordable:                 true,
				MonthlyShortage:              rub(0),
			},
		},
		{
			name: "marks goal as not affordable and calculates monthly shortage",
			goal: SavingsGoal{
				Name:          "Bike",
				TargetAmount:  rub(1200),
				CurrentAmount: rub(0),
			},
			months:        10,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			want: GoalDeadlineResult{
				TargetAmount:                 rub(1200),
				CurrentAmount:                rub(0),
				RemainingAmount:              rub(1200),
				Months:                       10,
				RequiredMonthlyContribution:  rub(120),
				AvailableMonthlyContribution: rub(100),
				IsAffordable:                 false,
				MonthlyShortage:              rub(20),
			},
		},
		{
			name: "handles partially completed goal",
			goal: SavingsGoal{
				Name:          "Course",
				TargetAmount:  rub(1000),
				CurrentAmount: rub(400),
			},
			months:        6,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			want: GoalDeadlineResult{
				TargetAmount:                 rub(1000),
				CurrentAmount:                rub(400),
				RemainingAmount:              rub(600),
				Months:                       6,
				RequiredMonthlyContribution:  rub(100),
				AvailableMonthlyContribution: rub(100),
				IsAffordable:                 true,
				MonthlyShortage:              rub(0),
			},
		},
		{
			name: "handles already reached goal",
			goal: SavingsGoal{
				Name:          "Trip",
				TargetAmount:  rub(500),
				CurrentAmount: rub(700),
			},
			months:        6,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(0)),
			want: GoalDeadlineResult{
				TargetAmount:                 rub(500),
				CurrentAmount:                rub(700),
				RemainingAmount:              rub(0),
				Months:                       6,
				RequiredMonthlyContribution:  rub(0),
				AvailableMonthlyContribution: rub(0),
				IsAffordable:                 true,
				MonthlyShortage:              rub(0),
			},
		},
		{
			name: "returns error for invalid months",
			goal: SavingsGoal{
				Name:          "Trip",
				TargetAmount:  rub(500),
				CurrentAmount: rub(0),
			},
			months:        0,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr:       ErrInvalidMonths,
		},
		{
			name: "returns error for currency mismatch",
			goal: SavingsGoal{
				Name:          "Trip",
				TargetAmount:  rub(500),
				CurrentAmount: rub(0),
			},
			months:        5,
			budgetSummary: budgetWithOptionalSpendingLimit(domain.NewEUR(100_00)),
			wantErr:       domain.ErrCurrencyMismatch,
		},
		{
			name: "returns error for zero target amount",
			goal: SavingsGoal{
				Name:          "Invalid",
				TargetAmount:  rub(0),
				CurrentAmount: rub(0),
			},
			months:        5,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr:       ErrInvalidGoalTarget,
		},
		{
			name: "returns error for negative target amount",
			goal: SavingsGoal{
				Name:          "Invalid",
				TargetAmount:  rub(-100),
				CurrentAmount: rub(0),
			},
			months:        5,
			budgetSummary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr:       ErrInvalidGoalTarget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateGoalDeadline(tt.goal, tt.months, tt.budgetSummary)
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

func budgetWithOptionalSpendingLimit(optionalSpendingLimit domain.Money) BudgetSummary {
	return BudgetSummary{OptionalSpendingLimit: optionalSpendingLimit}
}

func reachedOptions(targetAmount domain.Money, currentAmount domain.Money) []GoalPlanOption {
	remainingAmount := domain.ZeroMoney(targetAmount.Currency())
	monthlyContribution := domain.ZeroMoney(targetAmount.Currency())

	return []GoalPlanOption{
		{
			Type:                GoalPlanOptionTypeComfortable,
			TargetAmount:        targetAmount,
			CurrentAmount:       currentAmount,
			RemainingAmount:     remainingAmount,
			MonthlyContribution: monthlyContribution,
			Months:              0,
			IsAffordable:        true,
		},
		{
			Type:                GoalPlanOptionTypeBalanced,
			TargetAmount:        targetAmount,
			CurrentAmount:       currentAmount,
			RemainingAmount:     remainingAmount,
			MonthlyContribution: monthlyContribution,
			Months:              0,
			IsAffordable:        true,
		},
		{
			Type:                GoalPlanOptionTypeAggressive,
			TargetAmount:        targetAmount,
			CurrentAmount:       currentAmount,
			RemainingAmount:     remainingAmount,
			MonthlyContribution: monthlyContribution,
			Months:              0,
			IsAffordable:        true,
		},
	}
}
