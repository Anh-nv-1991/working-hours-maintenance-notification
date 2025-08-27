package repo

import (
	"context"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"
)

// PGReadingRepo implements ports.ReadingRepo using sqlc.
type PGReadingRepo struct {
	q *db.Queries
}

func NewReadingRepo(q *db.Queries) *PGReadingRepo { return &PGReadingRepo{q: q} }

var _ ports.ReadingRepo = (*PGReadingRepo)(nil)

func (r *PGReadingRepo) AddReading(ctx context.Context, arg db.AddReadingParams) (db.Reading, error) {
	return r.q.AddReading(ctx, arg)
}

func (r *PGReadingRepo) LastReading(ctx context.Context, deviceID int64) (db.Reading, error) {
	return r.q.LastReading(ctx, deviceID)
}
