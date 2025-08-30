package domain

import "time"

type Plan struct {
	ID            PlanID
	Name          string
	IntervalHours int
	Description   *string

	CreatedAt time.Time
	UpdatedAt time.Time
}
