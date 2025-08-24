4) Data Access (sqlc + pgx)
4.1 Cài deps & tools
# trong VS Code PowerShell
go env -w GOPRIVATE=
go get github.com/jackc/pgx/v5 github.com/jackc/pgx/v5/pgxpool
go get github.com/joho/godotenv
go get github.com/gin-gonic/gin

# tools (local dev)
# sqlc: https://docs.sqlc.dev/en/latest/overview/install.html (nên cài binary)
# migrate: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

4.2 Tạo cấu hình sqlc.yaml

sqlc.yaml (gốc repo):

version: "2"
sql:
  - schema: "db/migrations"
    queries: "db/queries"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "db/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: false
        output_db_file_name: "generated.go"

4.3 Viết queries (ngắn gọn nhưng đủ dùng)

Tạo folder db/queries/ với 4 file:

devices.sql

-- name: CreateDevice :one
INSERT INTO devices (id, name, metadata) VALUES ($1,$2,$3) RETURNING *;

-- name: GetDevice :one
SELECT * FROM devices WHERE id = $1;

-- name: ListDevices :many
SELECT * FROM devices ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: UpdateDevice :one
UPDATE devices SET name=$2, metadata=$3, updated_at=now() WHERE id=$1 RETURNING *;

-- name: DeleteDevice :exec
DELETE FROM devices WHERE id=$1;


readings.sql

-- name: AddReading :one
INSERT INTO readings (device_id, value, at) VALUES ($1,$2,$3) RETURNING *;

-- name: LastReadingByDevice :one
SELECT * FROM readings WHERE device_id=$1 ORDER BY at DESC LIMIT 1;


plans.sql

-- name: CreatePlan :one
INSERT INTO plans (id, device_id, threshold_hi, threshold_lo, window) 
VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: GetPlanForDevice :one
SELECT * FROM plans WHERE device_id=$1 ORDER BY created_at DESC LIMIT 1;


alerts.sql

-- name: CreateAlert :one
INSERT INTO alerts (id, device_id, reading_id, kind, message) 
VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: ListOpenAlerts :many
SELECT * FROM alerts WHERE serviced_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ResolveAlert :one
UPDATE alerts SET serviced_at=now(), serviced_note=$2 WHERE id=$1 RETURNING *;


Giả định ACE đã có migrations tạo 4 bảng devices, readings, plans, alerts với các cột như trên. Nếu cần mình in lại migration mẫu.

4.4 Generate code
sqlc generate


Sinh ra package db/sqlc với struct, phương thức typed.

4.5 Repo adapter (pgxpool) + domain

internal/domain/models.go

package domain

import "time"

type Device struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Metadata  []byte     `json:"metadata,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Reading struct {
	ID       int64     `json:"id"`
	DeviceID string    `json:"device_id"`
	Value    float64   `json:"value"`
	At       time.Time `json:"at"`
}

type Plan struct {
	ID          string  `json:"id"`
	DeviceID    string  `json:"device_id"`
	ThresholdHi *float64
	ThresholdLo *float64
	WindowSec   int32
}

type Alert struct {
	ID          string     `json:"id"`
	DeviceID    string     `json:"device_id"`
	ReadingID   int64      `json:"reading_id"`
	Kind        string     `json:"kind"`
	Message     string     `json:"message"`
	ServicedAt  *time.Time `json:"serviced_at,omitempty"`
	ServicedNote *string   `json:"serviced_note,omitempty"`
}


internal/repo/repo.go

package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"va-cv-ace/db/sqlc"
)

type Repo struct {
	q *db.Queries
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{q: db.New(pool)}
}
func (r *Repo) Q() *db.Queries { return r.q }

// Gợi ý: bọc thêm method domain-level nếu muốn decouple hẳn.

5) Ports (interfaces) để tách HTTP/DB

internal/ports/ports.go

package ports

import (
	"context"
	"time"
	"va-cv-ace/db/sqlc"
)

type DeviceRepo interface {
	CreateDevice(ctx context.Context, arg db.CreateDeviceParams) (db.Device, error)
	GetDevice(ctx context.Context, id string) (db.Device, error)
	ListDevices(ctx context.Context, limit, offset int32) ([]db.Device, error)
	UpdateDevice(ctx context.Context, arg db.UpdateDeviceParams) (db.Device, error)
	DeleteDevice(ctx context.Context, id string) error
}

type ReadingRepo interface {
	AddReading(ctx context.Context, arg db.AddReadingParams) (db.Reading, error)
	LastReadingByDevice(ctx context.Context, deviceID string) (db.Reading, error)
}

type PlanRepo interface {
	GetPlanForDevice(ctx context.Context, deviceID string) (db.Plan, error)
}

type AlertRepo interface {
	CreateAlert(ctx context.Context, arg db.CreateAlertParams) (db.Alert, error)
	ListOpenAlerts(ctx context.Context, limit, offset int32) ([]db.Alert, error)
	ResolveAlert(ctx context.Context, arg db.ResolveAlertParams) (db.Alert, error)
}

type Clock interface { Now() time.Time }


Nhanh gọn: vì sqlc đã tạo methods, ta có thể type-alias *db.Queries để implement các interface trên.

6) Use‑cases (AddReading, ComputeAlert, MarkServiced)

internal/usecase/add_reading.go

package usecase

import (
	"context"
	"fmt"
	"va-cv-ace/db/sqlc"
	"va-cv-ace/internal/ports"
)

type AddReadingUC struct {
	R  ports.ReadingRepo
	P  ports.PlanRepo
	A  ports.AlertRepo
}

func (uc *AddReadingUC) Exec(ctx context.Context, deviceID string, value float64) (db.Reading, *db.Alert, error) {
	rd, err := uc.R.AddReading(ctx, db.AddReadingParams{DeviceID: deviceID, Value: value, At: rdNow()})
	if err != nil { return db.Reading{}, nil, err }

	plan, err := uc.P.GetPlanForDevice(ctx, deviceID)
	if err != nil { return rd, nil, nil } // không có plan thì thôi

	var alert *db.Alert
	if plan.ThresholdHi != nil && value > *plan.ThresholdHi {
		a, err := uc.A.CreateAlert(ctx, db.CreateAlertParams{
			ID: RandID(), DeviceID: deviceID, ReadingID: rd.ID, Kind: "HI",
			Message: fmt.Sprintf("Value %.2f > hi %.2f", value, *plan.ThresholdHi),
		})
		if err == nil { alert = &a }
	}
	if plan.ThresholdLo != nil && value < *plan.ThresholdLo {
		a, err := uc.A.CreateAlert(ctx, db.CreateAlertParams{
			ID: RandID(), DeviceID: deviceID, ReadingID: rd.ID, Kind: "LO",
			Message: fmt.Sprintf("Value %.2f < lo %.2f", value, *plan.ThresholdLo),
		})
		if err == nil { alert = &a }
	}
	return rd, alert, nil
}


internal/usecase/mark_serviced.go

package usecase

import (
	"context"
	"va-cv-ace/db/sqlc"
	"va-cv-ace/internal/ports"
)

type MarkServicedUC struct{ A ports.AlertRepo }

func (uc *MarkServicedUC) Exec(ctx context.Context, id, note string) (db.Alert, error) {
	return uc.A.ResolveAlert(ctx, db.ResolveAlertParams{ID: id, ServicedNote: &note})
}


rdNow() và RandID() ACE có thể dùng time.Now().UTC() và github.com/google/uuid (go get) hoặc quick hack tự tạo.

7) HTTP API (Gin)

cmd/server/main.go (tối giản, chạy được):

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"

	"va-cv-ace/db/sqlc"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil { log.Fatal(err) }
	defer pool.Close()

	// cheap readiness ping
	if err := pingDB(ctx, pool); err != nil { log.Fatal(err) }

	q := db.New(pool)
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/readiness", func(c *gin.Context) {
		if err := pingDB(c, pool); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "err": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	// sample: list devices (ACE có thể thêm CRUD khác tương tự)
	r.GET("/devices", func(c *gin.Context) {
		rows, err := q.ListDevices(c, 50, 0)
		if err != nil { c.JSON(500, gin.H{"error": err.Error()}); return }
		c.JSON(200, rows)
	})

	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil { log.Fatal(err) }
}

func pingDB(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second); defer cancel()
	return pool.Ping(ctx)
}


configs/.env.example

APP_ENV=local
PORT=8080
DATABASE_URL=postgres://postgres:postgres@db:5432/app?sslmode=disable

8) Observability (nhẹ nhàng)

Gin đã có logger + recovery mặc định từ gin.Default(). CORS: thêm middleware sau nếu cần.

9) Dockerfile (multi‑stage, non‑root)

build/Dockerfile

# build
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/server

# run (distroless non-root)
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /home/nonroot
ENV PORT=8080
COPY --from=build /out/app /usr/local/bin/app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/app"]

10) docker‑compose (api + db + adminer optional)

docker-compose.yml

services:
  db:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: app
    ports: ["5432:5432"]
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d app"]
      interval: 3s
      timeout: 3s
      retries: 20

  migrate:
    image: migrate/migrate:latest
    depends_on: [db]
    volumes:
      - ./db/migrations:/migrations
    entrypoint: [
      "sh","-c",
      "migrate -path=/migrations -database postgres://postgres:postgres@db:5432/app?sslmode=disable up && sleep 2"
    ]
    restart: "on-failure"

  api:
    build:
      context: .
      dockerfile: build/Dockerfile
    env_file: [configs/.env.example]
    depends_on:
      db:
        condition: service_healthy
      migrate:
        condition: service_started
    ports: ["8080:8080"]

  adminer:
    image: adminer
    ports: ["8081:8080"]
    depends_on: [db]

volumes:
  pgdata:


Chạy thử:

docker compose up --build
# test thử:
Invoke-WebRequest http://localhost:8080/healthz
Invoke-WebRequest http://localhost:8080/readiness

11) CI (GitHub Actions) — vet, test, sqlc gen, migrate (apply & rollback)

.github/workflows/ci.yml

name: ci
on:
  pull_request:
  push:
    branches: [ main ]

jobs:
  build-test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: app
        ports: ["5432:5432"]
        options: >-
          --health-cmd="pg_isready -U postgres -d app"
          --health-interval=5s --health-timeout=5s --health-retries=10
    env:
      DATABASE_URL: postgres://postgres:postgres@localhost:5432/app?sslmode=disable
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Install sqlc
        run: |
          curl -sSL https://downloads.sqlc.dev/sqlc_1.28.0_linux_amd64.tar.gz | tar -xz
          sudo mv sqlc /usr/local/bin/sqlc
          sqlc version

      - name: Install migrate
        run: |
          curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar -xz
          sudo mv migrate /usr/local/bin/migrate
          migrate -version

      - name: sqlc generate
        run: sqlc generate

      - name: Migrate up (apply)
        run: migrate -database "$DATABASE_URL" -path db/migrations up

      - name: Go vet
        run: go vet ./...

      - name: Go test
        run: go test ./...

      - name: Migrate down (rollback)
        run: migrate -database "$DATABASE_URL" -path db/migrations down -all

  docker-build:
    runs-on: ubuntu-latest
    needs: build-test
    steps:
      - uses: actions/checkout@v4
      - name: Build Docker image
        run: docker build -f build/Dockerfile -t va-cv-ace:ci .


Note: Không có “dry‑run” thuần của migrate, nên mình apply lên Postgres service rồi rollback trong CI — an toàn vì DB ephemeral.

Quick checklist để xanh CI

db/migrations có đầy đủ schema cho 4 bảng và foreign key.

sqlc.yaml đúng đường dẫn.

db/queries/*.sql hợp lệ → sqlc generate pass.

go test ./... có thể để trống test ban đầu (tạo 1 test rỗng để pass).

internal/smoke/smoke_test.go

package smoke
import "testing"
func TestOK(t *testing.T) {}


Actions có thể tải sqlc + migrate thành công.

DATABASE_URL trong CI trỏ localhost (service postgres đã export port).