package port

import (
	"context"
	"time"
	"wh-ma/internal/domain"
)

type CreateMaintenanceInput struct {
	DeviceID    domain.DeviceID
	At          time.Time
	Interval    *int32 // nil nếu không áp dụng bậc interval
	Notes       *string
	PerformedBy *string
	Cost        *string // tiền tệ dạng decimal string, ví dụ "12345.67"
}

type MaintenanceRepository interface {
	Create(ctx context.Context, in CreateMaintenanceInput) (*domain.MaintenanceEvent, error)
	Delete(ctx context.Context, id int64) error
	ListByDevice(ctx context.Context, deviceID domain.DeviceID, limit, offset int32) ([]*domain.MaintenanceEvent, error)
}
