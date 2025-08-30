package request

import "time"

// POST /devices
type CreateDevice struct {
	SerialNumber   string     `json:"serial_number" binding:"required"`
	Name           string     `json:"name" binding:"required"`
	Model          string     `json:"model"`
	Manufacturer   string     `json:"manufacturer"`
	Year           int        `json:"year"`
	CommissionDate *time.Time `json:"commission_date"` // RFC3339
	Status         string     `json:"status"`          // optional; default "active"
	Location       *string    `json:"location"`
	PlanID         *int64     `json:"plan_id"`
}

// PATCH /devices/:id
type UpdateBasic struct {
	Name     string  `json:"name" binding:"required"`
	Status   string  `json:"status" binding:"required"`
	Location *string `json:"location"` // optional
}

// PATCH /devices/:id/plan  (nil -> bỏ plan)
type UpdatePlan struct {
	PlanID *int64 `json:"plan_id"`
}
