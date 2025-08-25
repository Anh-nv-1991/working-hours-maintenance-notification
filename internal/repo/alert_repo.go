package repo

import (
	"context"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"
)

// PGAlertRepo implements ports.AlertRepo using sqlc.
type PGAlertRepo struct {
	q *db.Queries
}

func NewAlertRepo(q *db.Queries) *PGAlertRepo { return &PGAlertRepo{q: q} }

var _ ports.AlertRepo = (*PGAlertRepo)(nil)

func (r *PGAlertRepo) CreateAlert(ctx context.Context, arg db.CreateAlertParams) (db.Alert, error) {
	return r.q.CreateAlert(ctx, arg)
}

func (r *PGAlertRepo) ListOpenAlertsByDevice(ctx context.Context, deviceID int64) ([]db.Alert, error) {
	return r.q.ListOpenAlertsByDevice(ctx, deviceID)
}

func (r *PGAlertRepo) ResolveAlert(ctx context.Context, arg db.ResolveAlertParams) (db.Alert, error) {
	return r.q.ResolveAlert(ctx, arg)
}

func (r *PGAlertRepo) CheckThresholdBreach(ctx context.Context, deviceID int64) (db.CheckThresholdBreachRow, error) {
	return r.q.CheckThresholdBreach(ctx, deviceID)
}
