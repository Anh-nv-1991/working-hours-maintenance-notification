package repo

import (
	"context"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"
)

// PGPlanRepo implements ports.PlanRepo using sqlc.
type PGPlanRepo struct {
	q *db.Queries
}

func NewPlanRepo(q *db.Queries) *PGPlanRepo { return &PGPlanRepo{q: q} }

var _ ports.PlanRepo = (*PGPlanRepo)(nil)

func (r *PGPlanRepo) GetPlanForDevice(ctx context.Context, deviceID int64) (db.Plan, error) {
	return r.q.GetPlanForDevice(ctx, deviceID)
}

func (r *PGPlanRepo) UpsertPlan(ctx context.Context, arg db.UpsertPlanParams) (db.Plan, error) {
	return r.q.UpsertPlan(ctx, arg)
}
