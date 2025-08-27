package ports

import (
	"context"
	"maint/internal/domain"
)

type Repository interface {
	// readings
	AddReading(ctx context.Context, deviceID int64, value float64, atUnix int64) (*domain.Reading, error)
	GetLastReading(ctx context.Context, deviceID int64) (*domain.Reading, error)

	// plans
	GetPlanByDevice(ctx context.Context, deviceID int64) (*domain.Plan, error)

	// alerts
	CreateAlert(ctx context.Context, deviceID int64, readingID *int64, message string) (*domain.Alert, error)
	GetOpenAlertByDevice(ctx context.Context, deviceID int64) (*domain.Alert, error)
	MarkAlertServiced(ctx context.Context, alertID int64, servicedAt int64) (*domain.Alert, error)

	// transaction helper
	WithTx(ctx context.Context, fn func(r Repository) error) error
}

type Clock interface {
	NowUnix() int64
}

type Notifier interface {
	NotifyAlert(ctx context.Context, a *domain.Alert) error
}
