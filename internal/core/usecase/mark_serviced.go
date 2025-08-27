// internal/usecase/mark_serviced.go
package usecase

import (
	"context"

	db "maint/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type MarkServicedInput struct {
	AlertID int64
	// Optional: nếu muốn set thời điểm cụ thể, thêm field *time.Time và map Valid=true
}

type MarkServicedUseCase struct {
	AlertsRepo interface {
		ResolveAlert(ctx context.Context, arg db.ResolveAlertParams) (db.Alert, error)
	}
}

func NewMarkServicedUseCase(
	alertsRepo interface {
		ResolveAlert(ctx context.Context, arg db.ResolveAlertParams) (db.Alert, error)
	},
) *MarkServicedUseCase {
	return &MarkServicedUseCase{AlertsRepo: alertsRepo}
}

func (uc *MarkServicedUseCase) Execute(ctx context.Context, in MarkServicedInput) error {
	// Sửa ở đây: dùng pgtype.Timestamptz{} thay cho db.NullTimestamptz
	_, err := uc.AlertsRepo.ResolveAlert(ctx, db.ResolveAlertParams{
		ID:         in.AlertID,
		ServicedAt: pgtype.Timestamptz{}, // Valid = false -> DB tự NOW() theo COALESCE($2, NOW()) :contentReference[oaicite:4]{index=4}
	})
	return err
}
