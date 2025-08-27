package usecase

import (
	"context"
	"maint/internal/ports"
	"time"
)

// Clock/Notifier để dễ test/mock
type Clock interface {
	Now() time.Time
}
type Notifier interface {
	NotifyAlert(ctx context.Context, deviceID int64, level string, msg string) error
}

type Deps struct {
	Devices  ports.DeviceRepo
	Readings ports.ReadingRepo
	Plans    ports.PlanRepo
	Alerts   ports.AlertRepo
	Clock    Clock    // optional
	Notifier Notifier // optional
}
