package domain

type FinancialPlan struct {
	User              User
	Accounts          []Account
	Debts             []Debt
	Goals             []Goal
	RecurringIncomes  []RecurringIncome
	RecurringExpenses []RecurringExpense
}
