package ai

import (
	"encoding/json"
	"fmt"

	"finance/internal/calculator"
)

type PurchaseExplanationPromptData struct {
	BudgetSummary    BudgetSummaryPromptData    `json:"budgetSummary"`
	PurchaseDecision PurchaseDecisionPromptData `json:"purchaseDecision"`
}

type BudgetSummaryPromptData struct {
	CurrentBalance        string `json:"currentBalance"`
	ExpectedIncome        string `json:"expectedIncome"`
	RequiredExpenses      string `json:"requiredExpenses"`
	DebtPayments          string `json:"debtPayments"`
	GoalContributions     string `json:"goalContributions"`
	FreeMoney             string `json:"freeMoney"`
	OptionalSpendingLimit string `json:"optionalSpendingLimit"`
}

type PurchaseDecisionPromptData struct {
	Type           string `json:"type"`
	PurchaseName   string `json:"purchaseName"`
	Price          string `json:"price"`
	AvailableMoney string `json:"availableMoney"`
	MoneyLeft      string `json:"moneyLeft"`
	Shortage       string `json:"shortage"`
	MonthsToSave   int    `json:"monthsToSave"`
}

func BuildPurchaseExplanationPrompt(
	summary calculator.BudgetSummary,
	decision calculator.PurchaseDecision,
) (string, error) {
	data := PurchaseExplanationPromptData{
		BudgetSummary: BudgetSummaryPromptData{
			CurrentBalance:        summary.CurrentBalance.String(),
			ExpectedIncome:        summary.ExpectedIncome.String(),
			RequiredExpenses:      summary.RequiredExpenses.String(),
			DebtPayments:          summary.DebtPayments.String(),
			GoalContributions:     summary.GoalContributions.String(),
			FreeMoney:             summary.FreeMoney.String(),
			OptionalSpendingLimit: summary.OptionalSpendingLimit.String(),
		},
		PurchaseDecision: PurchaseDecisionPromptData{
			Type:           string(decision.Type),
			PurchaseName:   decision.PurchaseName,
			Price:          decision.Price.String(),
			AvailableMoney: decision.AvailableMoney.String(),
			MoneyLeft:      decision.MoneyLeft.String(),
			Shortage:       decision.Shortage.String(),
			MonthsToSave:   decision.MonthsToSave,
		},
	}

	rawData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`
You are a friendly financial assistant.

You receive calculated financial data from a deterministic calculator.
The calculator is the source of truth.

Important rules:
- Do not recalculate numbers.
- Do not change the purchase decision.
- Do not invent missing income, expenses, debts, goals, risks, or user preferences.
- Use only the provided calculated values.
- Answer in Russian.
- Keep the answer short, warm, and practical.
- Do not use raw "RUB" formatting if possible. Rewrite money values in a human-friendly Russian format, for example "58 000 ₽".
- Do not mention that you are an AI or that you received JSON.

Decision rules:
- If purchaseDecision.type is "affordable":
  - Say that the purchase fits into the free monthly budget.
  - Say how much money will be left after the purchase.
  - Explain that required expenses, debt payments, and goal contributions are already accounted for.
- If purchaseDecision.type is "not_affordable":
  - Say that the purchase does not fit into the free monthly budget.
  - Say that it exceeds the available money by purchaseDecision.shortage.
  - Do not say that "money will be left" after the purchase.
  - Do not imply the user can buy it now.
  - If monthsToSave is greater than 0, explain that the safest option is to postpone the purchase for approximately that many months.
  - If goalContributions are greater than zero, you may mention that temporarily reducing goal contributions is a possible trade-off, but it will slow down the goal.
  - Do not recommend reducing specific expenses unless specific optional expenses are provided.

Answer structure:
1. Verdict: one direct sentence.
2. Explanation: why this decision was made.
3. Trade-off: what would need to change to buy it.
4. Recommendation: one practical next step.

Calculated data:
%s
`, string(rawData))

	return prompt, nil
}
