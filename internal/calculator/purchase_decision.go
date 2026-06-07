package calculator

import (
	"errors"
	"fmt"

	"finance/internal/domain"
)

var ErrInvalidPurchasePrice = errors.New("purchase price must be greater than zero")

type PurchaseRequest struct {
	Name  string
	Price domain.Money
}

type PurchaseDecisionType string

const (
	PurchaseDecisionTypeAffordable    PurchaseDecisionType = "affordable"
	PurchaseDecisionTypeNotAffordable PurchaseDecisionType = "not_affordable"
)

type PurchaseDecision struct {
	Type           PurchaseDecisionType
	PurchaseName   string
	Price          domain.Money
	AvailableMoney domain.Money
	MoneyLeft      domain.Money
	Shortage       domain.Money
	MonthsToSave   int
	Reasons        []string
	Suggestions    []string
}

func CalculatePurchaseDecision(purchase PurchaseRequest, summary BudgetSummary) (PurchaseDecision, error) {
	availableMoney := summary.OptionalSpendingLimit
	if err := validatePurchaseRequest(purchase, availableMoney); err != nil {
		return PurchaseDecision{}, err
	}

	zero := domain.ZeroMoney(purchase.Price.Currency())
	decision := PurchaseDecision{
		PurchaseName:   purchase.Name,
		Price:          purchase.Price,
		AvailableMoney: availableMoney,
		MoneyLeft:      zero,
		Shortage:       zero,
	}

	if purchase.Price.MinorUnits() <= availableMoney.MinorUnits() {
		moneyLeft, err := availableMoney.Sub(purchase.Price)
		if err != nil {
			return PurchaseDecision{}, err
		}

		decision.Type = PurchaseDecisionTypeAffordable
		decision.MoneyLeft = moneyLeft
		decision.Reasons = []string{
			fmt.Sprintf("%s fits into the free monthly budget.", purchase.Name),
		}
		decision.Suggestions = []string{
			"This purchase does not touch required expenses, debt payments, or goal contributions.",
		}

		return decision, nil
	}

	shortage, err := purchase.Price.Sub(availableMoney)
	if err != nil {
		return PurchaseDecision{}, err
	}

	decision.Type = PurchaseDecisionTypeNotAffordable
	decision.Shortage = shortage
	decision.Reasons = []string{
		fmt.Sprintf("%s exceeds the available monthly money.", purchase.Name),
	}

	if availableMoney.MinorUnits() > 0 {
		decision.MonthsToSave = ceilDivide(shortage.MinorUnits(), availableMoney.MinorUnits())
		decision.Suggestions = []string{
			fmt.Sprintf("Save the shortage for %d month(s) before buying this.", decision.MonthsToSave),
		}
	} else {
		decision.Suggestions = []string{
			"Increase income or reduce expenses before buying this.",
		}
	}

	return decision, nil
}

func validatePurchaseRequest(purchase PurchaseRequest, availableMoney domain.Money) error {
	if purchase.Price.MinorUnits() <= 0 {
		return ErrInvalidPurchasePrice
	}

	if purchase.Price.Currency() != availableMoney.Currency() {
		return domain.ErrCurrencyMismatch
	}

	return nil
}
