package domain

import "time"

type GoalStatus string

const (
	GoalStatusActive    GoalStatus = "active"
	GoalStatusCompleted GoalStatus = "completed"
	GoalStatusPaused    GoalStatus = "paused"
)

type Goal struct {
	ID            string
	UserID        string
	Name          string
	TargetAmount  Money
	CurrentAmount Money
	Deadline      time.Time
	Status        GoalStatus
}
