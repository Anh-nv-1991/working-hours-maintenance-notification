// internal/bootstrap/bootstrap.go
package bootstrap

import (
	"context"
	"os"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"

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
	q := db.New(pool)

	return &Deps{
		Devices:  q,
		Readings: q,
		Plans:    q,
		Alerts:   q,
		Pool:     pool,
	}, nil
}
