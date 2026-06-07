package calculator

import (
	"errors"

	"finance/internal/domain"
)

var (
	ErrInvalidGoalTarget = errors.New("goal target amount must be greater than zero")
	ErrInvalidMonths     = errors.New("months must be greater than zero")
)

type SavingsGoal struct {
	Name          string
	TargetAmount  domain.Money
	CurrentAmount domain.Money
}

type GoalPlanOptionType string

const (
	GoalPlanOptionTypeComfortable GoalPlanOptionType = "comfortable"
	GoalPlanOptionTypeBalanced    GoalPlanOptionType = "balanced"
	GoalPlanOptionTypeAggressive  GoalPlanOptionType = "aggressive"
)

type GoalPlanOption struct {
	Type                GoalPlanOptionType
	TargetAmount        domain.Money
	CurrentAmount       domain.Money
	RemainingAmount     domain.Money
	MonthlyContribution domain.Money
	Months              int
	IsAffordable        bool
}

type GoalDeadlineResult struct {
	TargetAmount                 domain.Money
	CurrentAmount                domain.Money
	RemainingAmount              domain.Money
	Months                       int
	RequiredMonthlyContribution  domain.Money
	AvailableMonthlyContribution domain.Money
	IsAffordable                 bool
	MonthlyShortage              domain.Money
}

func CalculateGoalPlanOptions(goal SavingsGoal, budgetSummary BudgetSummary) ([]GoalPlanOption, error) {
	availableMonthlyContribution := budgetSummary.OptionalSpendingLimit
	if err := validateSavingsGoal(goal, availableMonthlyContribution); err != nil {
		return nil, err
	}

	remainingAmount, err := calculateRemainingGoalAmount(goal)
	if err != nil {
		return nil, err
	}

	if remainingAmount.MinorUnits() == 0 {
		return reachedGoalPlanOptions(goal, remainingAmount), nil
	}

	if availableMonthlyContribution.MinorUnits() <= 0 {
		return []GoalPlanOption{}, nil
	}

	plans := []struct {
		optionType GoalPlanOptionType
		percent    int64
	}{
		{optionType: GoalPlanOptionTypeComfortable, percent: 50},
		{optionType: GoalPlanOptionTypeBalanced, percent: 80},
		{optionType: GoalPlanOptionTypeAggressive, percent: 100},
	}

	options := make([]GoalPlanOption, 0, len(plans))
	for _, plan := range plans {
		monthlyContribution := moneyPercent(availableMonthlyContribution, plan.percent)
		if monthlyContribution.MinorUnits() <= 0 {
			continue
		}

		options = append(options, GoalPlanOption{
			Type:                plan.optionType,
			TargetAmount:        goal.TargetAmount,
			CurrentAmount:       goal.CurrentAmount,
			RemainingAmount:     remainingAmount,
			MonthlyContribution: monthlyContribution,
			Months:              ceilDivide(remainingAmount.MinorUnits(), monthlyContribution.MinorUnits()),
			IsAffordable:        true,
		})
	}

	return options, nil
}

func CalculateGoalDeadline(goal SavingsGoal, months int, budgetSummary BudgetSummary) (GoalDeadlineResult, error) {
	if months <= 0 {
		return GoalDeadlineResult{}, ErrInvalidMonths
	}

	availableMonthlyContribution := budgetSummary.OptionalSpendingLimit
	if err := validateSavingsGoal(goal, availableMonthlyContribution); err != nil {
		return GoalDeadlineResult{}, err
	}

	remainingAmount, err := calculateRemainingGoalAmount(goal)
	if err != nil {
		return GoalDeadlineResult{}, err
	}

	zero := domain.ZeroMoney(goal.TargetAmount.Currency())
	result := GoalDeadlineResult{
		TargetAmount:                 goal.TargetAmount,
		CurrentAmount:                goal.CurrentAmount,
		RemainingAmount:              remainingAmount,
		Months:                       months,
		AvailableMonthlyContribution: availableMonthlyContribution,
		MonthlyShortage:              zero,
	}

	if remainingAmount.MinorUnits() == 0 {
		result.RequiredMonthlyContribution = zero
		result.IsAffordable = true
		return result, nil
	}

	requiredMonthlyContribution := domain.NewMoney(
		int64(ceilDivide(remainingAmount.MinorUnits(), int64(months))),
		goal.TargetAmount.Currency(),
	)
	result.RequiredMonthlyContribution = requiredMonthlyContribution
	result.IsAffordable = requiredMonthlyContribution.MinorUnits() <= availableMonthlyContribution.MinorUnits()

	if !result.IsAffordable {
		monthlyShortage, err := requiredMonthlyContribution.Sub(availableMonthlyContribution)
		if err != nil {
			return GoalDeadlineResult{}, err
		}

		result.MonthlyShortage = monthlyShortage
	}

	return result, nil
}

func validateSavingsGoal(goal SavingsGoal, availableMonthlyContribution domain.Money) error {
	if goal.TargetAmount.MinorUnits() <= 0 {
		return ErrInvalidGoalTarget
	}

	if goal.TargetAmount.Currency() != goal.CurrentAmount.Currency() ||
		goal.TargetAmount.Currency() != availableMonthlyContribution.Currency() {
		return domain.ErrCurrencyMismatch
	}

	return nil
}

func calculateRemainingGoalAmount(goal SavingsGoal) (domain.Money, error) {
	if goal.CurrentAmount.MinorUnits() >= goal.TargetAmount.MinorUnits() {
		return domain.ZeroMoney(goal.TargetAmount.Currency()), nil
	}

	return goal.TargetAmount.Sub(goal.CurrentAmount)
}

func reachedGoalPlanOptions(goal SavingsGoal, remainingAmount domain.Money) []GoalPlanOption {
	zero := domain.ZeroMoney(goal.TargetAmount.Currency())
	optionTypes := []GoalPlanOptionType{
		GoalPlanOptionTypeComfortable,
		GoalPlanOptionTypeBalanced,
		GoalPlanOptionTypeAggressive,
	}

	options := make([]GoalPlanOption, 0, len(optionTypes))
	for _, optionType := range optionTypes {
		options = append(options, GoalPlanOption{
			Type:                optionType,
			TargetAmount:        goal.TargetAmount,
			CurrentAmount:       goal.CurrentAmount,
			RemainingAmount:     remainingAmount,
			MonthlyContribution: zero,
			Months:              0,
			IsAffordable:        true,
		})
	}

	return options
}

func moneyPercent(money domain.Money, percent int64) domain.Money {
	return domain.NewMoney(money.MinorUnits()*percent/100, money.Currency())
}

func ceilDivide(amount int64, divisor int64) int {
	return int((amount + divisor - 1) / divisor)
}
