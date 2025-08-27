package repo

import (
	"context"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"
)

// PGDeviceRepo implements ports.DeviceRepo using sqlc.
type PGDeviceRepo struct {
	q *db.Queries
}

func NewDeviceRepo(q *db.Queries) *PGDeviceRepo { return &PGDeviceRepo{q: q} }

// compile-time check
var _ ports.DeviceRepo = (*PGDeviceRepo)(nil)

func (r *PGDeviceRepo) CreateDevice(ctx context.Context, name string) (db.Device, error) {
	return r.q.CreateDevice(ctx, name)
}

func (r *PGDeviceRepo) GetDevice(ctx context.Context, id int64) (db.Device, error) {
	return r.q.GetDevice(ctx, id)
}
