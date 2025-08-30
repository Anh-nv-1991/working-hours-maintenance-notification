package port

import (
	"context"

	"wh-ma/internal/domain"
)

type CreateAlertInput struct {
	DeviceID domain.DeviceID
	Type     string
	Message  string
}

type ResolveAlertInput struct {
	ID         int64
	ResolvedBy *string
}

type AlertRepository interface {
	Create(ctx context.Context, in CreateAlertInput) (*domain.Alert, error)
	ListOpenByDevice(ctx context.Context, deviceID domain.DeviceID, limit, offset int32) ([]*domain.Alert, error)
	Resolve(ctx context.Context, in ResolveAlertInput) (*domain.Alert, error)
}
