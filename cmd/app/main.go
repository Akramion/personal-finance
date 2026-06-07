package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"finance/ai"
	"finance/internal/calculator"
	"finance/internal/domain"

	"github.com/joho/godotenv"
)

func main() {
	user := domain.User{
		ID:    "user-1",
		Name:  "manual-budget",
		Email: "manual-budget@example.com",
	}

	account := domain.Account{
		ID:      "account-1",
		UserID:  user.ID,
		Name:    "Current account",
		Type:    domain.AccountTypeCard,
		Balance: rub(0),
	}

	salary := domain.RecurringIncome{
		ID:          "income-1",
		UserID:      user.ID,
		AccountID:   account.ID,
		Amount:      rub(450_000),
		Type:        domain.IncomeTypeSalary,
		Description: "Salary",
		PaymentDay:  1,
	}

	requiredExpenses := []domain.RecurringExpense{
		requiredExpense("expense-1", user.ID, account.ID, rub(47_000), domain.ExpenseCategoryRent, "Rent / Apartment", 1),
		requiredExpense("expense-2", user.ID, account.ID, rub(12_000), domain.ExpenseCategoryOther, "Piano lessons", 1),
		requiredExpense("expense-3", user.ID, account.ID, rub(24_000), domain.ExpenseCategoryOther, "English lessons", 1),
		requiredExpense("expense-4", user.ID, account.ID, rub(10_000), domain.ExpenseCategoryOther, "Driving lessons", 1),
		requiredExpense("expense-5", user.ID, account.ID, rub(2_000), domain.ExpenseCategoryOther, "Phone", 1),
		requiredExpense("expense-6", user.ID, account.ID, rub(20_000), domain.ExpenseCategoryOther, "Yandex Split", 1),
		requiredExpense("expense-7", user.ID, account.ID, rub(0), domain.ExpenseCategoryOther, "Smirnov School", 1),
		requiredExpense("expense-8", user.ID, account.ID, rub(4_000), domain.ExpenseCategoryOther, "Subscriptions", 1),
		requiredExpense("expense-9", user.ID, account.ID, rub(14_000), domain.ExpenseCategoryOther, "Dog/Cat", 1),
		requiredExpense("expense-10", user.ID, account.ID, rub(50_000), domain.ExpenseCategoryOther, "Food", 1),
	}

	goalContributions := rub(200_000)
	basePlan := domain.FinancialPlan{
		User:              user,
		Accounts:          []domain.Account{account},
		RecurringIncomes:  []domain.RecurringIncome{salary},
		RecurringExpenses: requiredExpenses,
	}

	planWithDebt := basePlan
	planWithDebt.Debts = []domain.Debt{
		{
			ID:              "debt-1",
			UserID:          user.ID,
			Name:            "Manual monthly debt payment",
			Type:            domain.DebtTypeOther,
			TotalAmount:     rub(45_000),
			RemainingAmount: rub(45_000),
			MonthlyPayment:  rub(45_000),
			PaymentDay:      1,
			CreatedAt:       time.Now(),
		},
	}

	printBudgetSummary("Monthly budget with 45,000 RUB debt payment", planWithDebt, goalContributions)

	goal := calculator.SavingsGoal{
		Name:          "Переезд в Сербию",
		TargetAmount:  rub(1_000_000),
		CurrentAmount: rub(360_000),
	}

	summary, err := calculator.CalculateBudgetSummary(planWithDebt, goalContributions)
	if err != nil {
		return
	}

	fmt.Println(calculator.CalculateGoalPlanOptions(goal, summary))

	purchase := calculator.PurchaseRequest{
		Name:  "Akai MPC one+",
		Price: rub(80_000),
	}

	decision, err := calculator.CalculatePurchaseDecision(purchase, summary)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	if err := godotenv.Load(); err != nil {
		log.Println(".env not loaded:", err)
	}

	ctx := context.Background()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is empty")
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	llmClient, err := ai.NewOpenAIClient(apiKey, model)
	if err != nil {
		log.Fatal(err)
	}

	explanationService := ai.NewPurchaseExplanationService(llmClient)

	explanation, err := explanationService.Explain(ctx, summary, decision)
	if err != nil {
		fmt.Println("AI explanation failed, fallback was used:")
		fmt.Println("Error:", err)
	}

	fmt.Println(explanation)
}

func requiredExpense(
	id string,
	userID string,
	accountID string,
	amount domain.Money,
	category domain.ExpenseCategory,
	description string,
	paymentDay int,
) domain.RecurringExpense {
	return domain.RecurringExpense{
		ID:          id,
		UserID:      userID,
		AccountID:   accountID,
		Amount:      amount,
		Type:        domain.ExpenseTypeRequired,
		Category:    category,
		Description: description,
		PaymentDay:  paymentDay,
	}
}

func printBudgetSummary(title string, plan domain.FinancialPlan, goalContributions domain.Money) {
	summary, err := calculator.CalculateBudgetSummary(plan, goalContributions)
	if err != nil {
		fmt.Println(title+":", err)
		return
	}

	fmt.Println(title)
	fmt.Println("Current balance:", summary.CurrentBalance)
	fmt.Println("Expected income:", summary.ExpectedIncome)
	fmt.Println("Required expenses:", summary.RequiredExpenses)
	fmt.Println("Debt payments:", summary.DebtPayments)
	fmt.Println("Goal contributions:", summary.GoalContributions)
	fmt.Println("Free money:", summary.FreeMoney)
	fmt.Println("Optional spending limit:", summary.OptionalSpendingLimit)
	fmt.Println()
}

func rub(majorUnits int64) domain.Money {
	return domain.NewMoneyFromMajorUnits(majorUnits, domain.CurrencyRUB)
}
