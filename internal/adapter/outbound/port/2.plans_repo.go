package port

import (
	"context"
	"wh-ma/internal/domain"
)

type CreatePlanInput struct {
	Name          string
	IntervalHours int
	Description   *string
}

type UpdatePlanInput struct {
	ID            domain.PlanID
	Name          string
	IntervalHours int
	Description   *string
}

type PlanRepository interface {
	Create(ctx context.Context, in CreatePlanInput) (*domain.Plan, error)
	GetByID(ctx context.Context, id domain.PlanID) (*domain.Plan, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.Plan, error)
	Update(ctx context.Context, in UpdatePlanInput) (*domain.Plan, error)
	Delete(ctx context.Context, id domain.PlanID) error
}
