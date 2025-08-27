package ports

import (
	"context"
	db "maint/internal/db/sqlc"
)

/* DEVICE — sqlc có 2 hàm: CreateDevice(name) và GetDevice(id) */
type DeviceRepo interface {
	CreateDevice(ctx context.Context, name string) (db.Device, error) // KHÔNG có CreateDeviceParams
	GetDevice(ctx context.Context, id int64) (db.Device, error)       // id là int64
	// Nếu muốn List/Update/Delete => phải thêm query rồi sqlc generate
}

/* READING */
type ReadingRepo interface {
	AddReading(ctx context.Context, arg db.AddReadingParams) (db.Reading, error)
	LastReading(ctx context.Context, deviceID int64) (db.Reading, error)
}

/* PLAN */
type PlanRepo interface {
	GetPlanForDevice(ctx context.Context, deviceID int64) (db.Plan, error)
	UpsertPlan(ctx context.Context, arg db.UpsertPlanParams) (db.Plan, error)
}

/* ALERT */
type AlertRepo interface {
	CreateAlert(ctx context.Context, arg db.CreateAlertParams) (db.Alert, error)
	ListOpenAlertsByDevice(ctx context.Context, deviceID int64) ([]db.Alert, error)
	ResolveAlert(ctx context.Context, arg db.ResolveAlertParams) (db.Alert, error)
	CheckThresholdBreach(ctx context.Context, deviceID int64) (db.CheckThresholdBreachRow, error)
}
