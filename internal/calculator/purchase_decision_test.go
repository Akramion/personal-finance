package calculator

import (
	"errors"
	"testing"

	"finance/internal/domain"
)

func TestCalculatePurchaseDecision(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		purchase       PurchaseRequest
		summary        BudgetSummary
		want           PurchaseDecision
		wantErr        error
		wantReason     string
		wantSuggestion string
	}{
		{
			name: "affordable purchase",
			purchase: PurchaseRequest{
				Name:  "Headphones",
				Price: rub(40),
			},
			summary: budgetWithOptionalSpendingLimit(rub(100)),
			want: PurchaseDecision{
				Type:           PurchaseDecisionTypeAffordable,
				PurchaseName:   "Headphones",
				Price:          rub(40),
				AvailableMoney: rub(100),
				MoneyLeft:      rub(60),
				Shortage:       rub(0),
				MonthsToSave:   0,
			},
			wantReason:     "Headphones fits into the free monthly budget.",
			wantSuggestion: "This purchase does not touch required expenses, debt payments, or goal contributions.",
		},
		{
			name: "not affordable purchase",
			purchase: PurchaseRequest{
				Name:  "Laptop",
				Price: rub(250),
			},
			summary: budgetWithOptionalSpendingLimit(rub(100)),
			want: PurchaseDecision{
				Type:           PurchaseDecisionTypeNotAffordable,
				PurchaseName:   "Laptop",
				Price:          rub(250),
				AvailableMoney: rub(100),
				MoneyLeft:      rub(0),
				Shortage:       rub(150),
				MonthsToSave:   3,
			},
			wantReason:     "Laptop exceeds the available monthly money.",
			wantSuggestion: "Save for 3 month(s) before buying this.",
		},
		{
			name: "months to save rounded up",
			purchase: PurchaseRequest{
				Name:  "Course",
				Price: rub(101),
			},
			summary: budgetWithOptionalSpendingLimit(rub(50)),
			want: PurchaseDecision{
				Type:           PurchaseDecisionTypeNotAffordable,
				PurchaseName:   "Course",
				Price:          rub(101),
				AvailableMoney: rub(50),
				MoneyLeft:      rub(0),
				Shortage:       rub(51),
				MonthsToSave:   3,
			},
			wantReason:     "Course exceeds the available monthly money.",
			wantSuggestion: "Save for 3 month(s) before buying this.",
		},
		{
			name: "zero available money",
			purchase: PurchaseRequest{
				Name:  "Bike",
				Price: rub(100),
			},
			summary: budgetWithOptionalSpendingLimit(rub(0)),
			want: PurchaseDecision{
				Type:           PurchaseDecisionTypeNotAffordable,
				PurchaseName:   "Bike",
				Price:          rub(100),
				AvailableMoney: rub(0),
				MoneyLeft:      rub(0),
				Shortage:       rub(100),
				MonthsToSave:   0,
			},
			wantReason:     "Bike exceeds the available monthly money.",
			wantSuggestion: "Increase income or reduce expenses before buying this.",
		},
		{
			name: "invalid purchase price",
			purchase: PurchaseRequest{
				Name:  "Invalid",
				Price: rub(0),
			},
			summary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr: ErrInvalidPurchasePrice,
		},
		{
			name: "negative purchase price",
			purchase: PurchaseRequest{
				Name:  "Invalid",
				Price: rub(-10),
			},
			summary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr: ErrInvalidPurchasePrice,
		},
		{
			name: "currency mismatch",
			purchase: PurchaseRequest{
				Name:  "Camera",
				Price: domain.NewUSD(100_00),
			},
			summary: budgetWithOptionalSpendingLimit(rub(100)),
			wantErr: domain.ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePurchaseDecision(tt.purchase, tt.summary)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Type != tt.want.Type ||
				got.PurchaseName != tt.want.PurchaseName ||
				got.Price != tt.want.Price ||
				got.AvailableMoney != tt.want.AvailableMoney ||
				got.MoneyLeft != tt.want.MoneyLeft ||
				got.Shortage != tt.want.Shortage ||
				got.MonthsToSave != tt.want.MonthsToSave {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}

			assertSingleMessage(t, "reason", got.Reasons, tt.wantReason)
			assertSingleMessage(t, "suggestion", got.Suggestions, tt.wantSuggestion)
		})
	}
}

func assertSingleMessage(t *testing.T, label string, got []string, want string) {
	t.Helper()

	if len(got) != 1 {
		t.Fatalf("expected one %s, got %d: %+v", label, len(got), got)
	}

	if got[0] != want {
		t.Fatalf("expected %s %q, got %q", label, want, got[0])
	}
}
