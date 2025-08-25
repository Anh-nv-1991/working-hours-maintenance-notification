// internal/bootstrap/bootstrap.go
package bootstrap

import (
	"context"
	"os"

	db "maint/internal/db/sqlc"
	"maint/internal/ports" // <— thêm
	"maint/internal/repo"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Deps struct {
	Devices  ports.DeviceRepo
	Readings ports.ReadingRepo
	Plans    ports.PlanRepo
	Alerts   ports.AlertRepo
	Pool     *pgxpool.Pool
}

func Init(ctx context.Context) (*Deps, error) {
	dsn := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	q := db.New(pool) // nếu sqlc.New cần pool hoặc conn

	devRepo := repo.NewDeviceRepo(q) // implements ports.DeviceRepo
	rdRepo := repo.NewReadingRepo(q) // implements ports.ReadingRepo
	plRepo := repo.NewPlanRepo(q)    // implements ports.PlanRepo
	alRepo := repo.NewAlertRepo(q)   // implements ports.AlertRepo

	return &Deps{
		Devices:  devRepo,
		Readings: rdRepo,
		Plans:    plRepo,
		Alerts:   alRepo,
		Pool:     pool,
	}, nil
}
