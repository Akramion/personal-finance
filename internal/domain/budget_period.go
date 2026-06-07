package domain

import "time"

type BudgetPeriod struct {
	ID        string
	UserID    string
	StartDate time.Time
	EndDate   time.Time
}
