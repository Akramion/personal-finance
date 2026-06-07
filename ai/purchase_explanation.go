package ai

import (
	"context"
	"fmt"
	"strings"

	"finance/internal/calculator"
)

type PurchaseExplanationService struct {
	client LLMClient
}

func NewPurchaseExplanationService(client LLMClient) PurchaseExplanationService {
	return PurchaseExplanationService{
		client: client,
	}
}

func (s PurchaseExplanationService) Explain(
	ctx context.Context,
	summary calculator.BudgetSummary,
	decision calculator.PurchaseDecision,
) (string, error) {
	prompt, err := BuildPurchaseExplanationPrompt(summary, decision)
	if err != nil {
		return "", err
	}

	explanation, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return BuildPurchaseExplanationFallback(decision), err
	}

	explanation = strings.TrimSpace(explanation)
	if explanation == "" {
		return BuildPurchaseExplanationFallback(decision), nil
	}

	return explanation, nil
}

func BuildPurchaseExplanationFallback(decision calculator.PurchaseDecision) string {
	switch decision.Type {
	case calculator.PurchaseDecisionTypeAffordable:
		return "Покупка помещается в свободный бюджет. После неё останется " +
			decision.MoneyLeft.String() + ". " +
			"Она не затрагивает обязательные расходы, долги и цели."

	case calculator.PurchaseDecisionTypeNotAffordable:
		if decision.MonthsToSave > 0 {
			return "Покупка пока не помещается в свободный бюджет. Не хватает " +
				decision.Shortage.String() + ". " +
				"Можно отложить недостающую сумму примерно за " +
				fmt.Sprintf("%d месяц(ев).", decision.MonthsToSave)
		}

		return "Покупка пока не помещается в свободный бюджет. Сейчас свободных денег нет, поэтому лучше сначала увеличить доход или снизить расходы."

	default:
		return "Не удалось подготовить подробное объяснение по покупке."
	}
}
