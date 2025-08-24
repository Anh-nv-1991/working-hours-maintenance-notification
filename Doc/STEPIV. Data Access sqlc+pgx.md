# Step 4 — Data Access with sqlc + pgx (Notes)

**Goal:** Turn hand‑written SQL into typed Go code, then call it via `pgxpool`.

**Terminal:** VS Code → **PowerShell**.

---

## TL;DR flow

```
db/migrations  +  db/queries/*.sql  --(sqlc generate)-->  internal/db/sqlc/*.go
```

* `sqlc` **does not touch DB**. It only reads schema + queries and generates Go code.
* Real DB calls happen when your Go app calls the generated functions with a live `DATABASE_URL`.

---

## 1) Create sqlc config (at repo root)

**File:** `sqlc.yaml`

```yaml
version: "2"
sql:
  - engine: postgresql
    schema: db/migrations
    queries: db/queries
    gen:
      go:
        package: db
        out: internal/db/sqlc
        sql_package: pgx/v5
```

> Adjust `schema:` if your migrations folder is different.

---

## 2) Queries we added (devices)

**Folder:** `db/queries`

**File:** `db/queries/devices.sql`

```sql
-- name: CreateDevice :one
INSERT INTO devices (name) VALUES ($1) RETURNING *;

-- name: GetDevice :one
SELECT * FROM devices WHERE id = $1;
```

> We removed `description` because current `devices` table doesn’t have that column.

---

## 3) Ensure schema (migrations)

Your migrations must define the `devices` table (Step 3). Minimal example if needed:

**File:** `db/migrations/000001_init.up.sql`

```sql
CREATE TABLE IF NOT EXISTS devices (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 4) Generate code using Docker (no local install needed)

Run at repo root:

```powershell
docker run --rm -v "${PWD}:/src" -w /src sqlc/sqlc generate
```

Output will appear in `internal/db/sqlc/` (e.g. `devices.sql.go`).

---

## 5) How the generated code is used (runtime)

* Create `pgxpool` with `DATABASE_URL` (e.g. `postgres://postgres:postgres@localhost:5432/app?sslmode=disable`).
* Build query handle: `q := db.New(pool)`.
* Call: `q.CreateDevice(ctx, "sensor-1")`, `q.GetDevice(ctx, id)`.

> The `-- name: <Func> :one|:many|:exec` comment controls the generated function name and return type.

---

## 6) Common errors & quick fixes

* **error: column "description" does not exist**

  * Fix query to match existing schema *or* add a migration to add the column, then re‑generate.
* **sqlc.(yaml|json) does not exist**

  * Create `sqlc.yaml` at repo root as above.
* **Couldn’t install sqlc on Windows**

  * Use the Docker command above (works everywhere).

---

## 7) Next (still Step 4)

Add more queries before generating again:

* `readings.sql`: `AddReading`, `LastReading`.
* `plans.sql`: `UpsertPlan`, `GetPlanForDevice`.
* `alerts.sql`: `CreateAlert`, `ResolveAlert`, `ListOpenAlertsByDevice`, `CheckThresholdBreach`.

Re‑run:

```powershell
docker run --rm -v "${PWD}:/src" -w /src sqlc/sqlc generate
```

---

## 8) Mini smoke test (optional)

Create `cmd/smoke/main.go` to call `CreateDevice` + `GetDevice` and print results. Ensure `DATABASE_URL` points to your running Postgres in Docker.
