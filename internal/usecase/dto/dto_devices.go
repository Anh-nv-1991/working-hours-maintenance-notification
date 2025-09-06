package dto

import (
	"time"
	"wh-ma/internal/domain"
)

type CreateDeviceCmd struct {
	SerialNumber   string
	Name           string
	Model          string
	Manufacturer   string
	Year           int
	CommissionDate *time.Time
	Status         domain.DeviceStatus
	Location       *string
	PlanID         *domain.PlanID
}

type UpdateDeviceBasicCmd struct {
	ID       domain.DeviceID
	Name     string
	Status   domain.DeviceStatus
	Location *string
}

type UpdateDevicePlanCmd struct {
	ID     domain.DeviceID
	PlanID *domain.PlanID // nil = bỏ kế hoạch
}
