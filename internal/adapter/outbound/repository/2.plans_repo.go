package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"wh-ma/internal/adapter/outbound/port"
	dbsqlc "wh-ma/internal/adapter/outbound/repository/sqlc"
	"wh-ma/internal/domain"
)

type PlanRepositoryPG struct {
	q *dbsqlc.Queries
}

func NewPlanRepository(pool *pgxpool.Pool) *PlanRepositoryPG {
	return &PlanRepositoryPG{q: dbsqlc.New(pool)}
}

// compile-time check
var _ port.PlanRepository = (*PlanRepositoryPG)(nil)

func (r *PlanRepositoryPG) Create(ctx context.Context, in port.CreatePlanInput) (*domain.Plan, error) {
	row, err := r.q.CreatePlan(ctx, dbsqlc.CreatePlanParams{
		Name:          in.Name,
		IntervalHours: int32(in.IntervalHours),
		Description:   in.Description, // *string
	})
	if err != nil {
		return nil, err
	}
	p := mapSqlcPlanToDomain(row)
	return &p, nil
}

func (r *PlanRepositoryPG) GetByID(ctx context.Context, id domain.PlanID) (*domain.Plan, error) {
	row, err := r.q.GetPlan(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	p := mapSqlcPlanToDomain(row)
	return &p, nil
}

func (r *PlanRepositoryPG) List(ctx context.Context, limit, offset int32) ([]*domain.Plan, error) {
	rows, err := r.q.ListPlans(ctx, dbsqlc.ListPlansParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Plan, 0, len(rows))
	for _, row := range rows {
		p := mapSqlcPlanToDomain(row)
		out = append(out, &p)
	}
	return out, nil
}

func (r *PlanRepositoryPG) Update(ctx context.Context, in port.UpdatePlanInput) (*domain.Plan, error) {
	row, err := r.q.UpdatePlan(ctx, dbsqlc.UpdatePlanParams{
		ID:            int64(in.ID),
		Name:          in.Name,
		IntervalHours: int32(in.IntervalHours),
		Description:   in.Description,
	})
	if err != nil {
		return nil, err
	}
	p := mapSqlcPlanToDomain(row)
	return &p, nil
}

func (r *PlanRepositoryPG) Delete(ctx context.Context, id domain.PlanID) error {
	return r.q.DeletePlan(ctx, int64(id))
}

// ===== mapping =====
func mapSqlcPlanToDomain(x dbsqlc.Plan) domain.Plan {
	return domain.Plan{
		ID:            domain.PlanID(x.ID),
		Name:          x.Name,
		IntervalHours: int(x.IntervalHours),
		Description:   x.Description,
		CreatedAt:     x.CreatedAt.Time,
		UpdatedAt:     x.UpdatedAt.Time,
	}
}
