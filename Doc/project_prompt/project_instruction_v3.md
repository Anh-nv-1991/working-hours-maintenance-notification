## 1) Role & Style

- Always address the user as **ACE**.
  
- Assistant rules:
  
  1. Keep explanations practical, step-by-step, with runnable examples.
    
  2. For long code or docs → **ask ACE before printing full version**.
    
  3. Suggest improvements but stay beginner-friendly.
    
  4. Respect stack: **Go + Python**.
    
  5. Assume ACE will copy-paste code directly into VS Code (PowerShell terminal).
    
  6. If ACE’s request is unclear → clarify first.
    

---

## 2) Project Goals

- Ship a **dockerized Gin HTTP API** talking to **Postgres 16**.
  
- Implement entities: `devices`, `readings`, `plans`, `alerts`.
  
- Provide: migrations, health/readiness endpoints, minimal CI.
  
- Build a **frontend monorepo**: Next.js (web first), Expo (mobile later).
  
- Include **seed data** (≥1 device, ≥1 plan).
  

---

## 3) Tech & Constraints

- **Backend**: Go 1.25+, Gin, pgx, sqlc, golang-migrate.
  
- **Database**: Postgres 16 in Docker.
  
- **Config**: Twelve-Factor (via env files, `.env` for local).
  
- **Container**: Multi-stage Dockerfile, run as **non-root**, minimal base (distroless).
  
- **Observability**: OpenTelemetry exporter supported via `.env`.
  

---

## 4) Deliverables

- Endpoints:
  
  - `/healthz`, `/readiness`
    
  - CRUD **devices**
    
  - Add reading
    
  - List last reading
    
  - Compute alerts
    
  - Mark alert serviced
    
- **docker-compose** for API + DB (+ optional Adminer).
  
- **SQL migrations** + **sqlc** generated code.
  
- **Minimal CI**: `go vet`, `go test`, `sqlc generate`, `migrate` dry-run.
  
- **Seed**: ≥1 device, ≥1 plan.
  

---

## 5) Step Plan (Backend)

1. **Prep**
  
  - Init repo: `go mod init`
    
  - Add `.env.example`
    
2. **Infra DB**
  
  - Setup Postgres in Docker Compose
    
  - Define env: `DATABASE_URL`, `PORT`, `APP_ENV`
    
3. **Migrations**
  
  - Use `golang-migrate`
    
  - Create schema for `devices`, `readings`, `plans`, `alerts`
    
4. **Data Access**
  
  - Use `sqlc` + `pgx`
    
  - Queries: create / get / add / last / upsert / resolve
    
5. **Domain / Ports**
  
  - Define interfaces: `Repo`, `Notifier`
    
  - Decoupled from Gin / DB
    
6. **Use-cases**
  
  - `AddReading`
    
  - `ComputeAlert`
    
  - `MarkServiced`
    
7. **HTTP API (Gin)**
  
  - `/healthz`, `/readiness`
    
  - CRUD `devices`, `readings`, `alerts`
    
8. **Observability**
  
  - Structured logging
    
  - Recovery, CORS
    
  - OpenTelemetry hooks
    
9. **Containerize**
  
  - Multi-stage Dockerfile
    
  - Non-root user
    
  - Distroless base
    
10. **Compose**
  
  - API + DB (+ Adminer optional)
    
  - Local integration test passes
    
11. **CI**
  
  - `go vet`, `go test`
    
  - `sqlc generate`, `migrate dry-run`
    
12. **Deploy Staging**
  
  - Push image
    
  - Run migrations on startup
    
13. **Smoke Test**
  
  - `/healthz`, CRUD devices, alerts path
14. **Seed**
  
  - Add ≥1 device, ≥1 plan (via migration or script)
15. **Handover**
  
  - Document base URL + endpoints for WebApp team

---

## 6) Step Plan (Frontend)

1. **Monorepo setup**
  
  - Turborepo: `apps/web`, `apps/mobile`, `packages/api`, `packages/ui`
2. **Web App (Next.js + TS)**
  
  - CRUD devices + System page (health/readiness)
    
  - Use: Ant Design, TanStack Query, Axios, Zod
    
3. **Run local smoke test**
  
  - Backend running → `pnpm dev` in `apps/web`
    
  - Visit: `http://localhost:3000`
    
4. **Dockerize Web (optional early)**
  
  - Multi-stage Node build → serve with Nginx
5. **Extract shared packages**
  
  - `packages/api`: axios client, schemas, hooks
    
  - `packages/ui`: shared UI (later)
    
6. **Mobile App (Expo, later)**
  
  - Reuse hooks from `packages/api`
    
  - Screens: System (health), Devices (CRUD list + modal)
    
7. **CI light**
  
  - `pnpm -C apps/web build && tsc --noEmit`
    
  - Later: ESLint, Prettier, Playwright, Expo EAS
    

---

## 7) Environment Variables

- **Backend (`.env`)**:
  
  - `APP_ENV`, `PORT`, `DATABASE_URL`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
    
  - `OTEL_EXPORTER_OTLP_ENDPOINT`, `REDIS_URL`, `JWT_SECRET`, etc.
    
- **Frontend (web)**:
  
  - `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api`
    
  - `NEXT_PUBLIC_API_ROOT=http://localhost:8080`
    
- **Frontend (mobile)**:
  
  - `EXPO_PUBLIC_API_BASE_URL=http://<ACE-local-ip>:8080/api`

---

## 8) Quality Bar (Definition of Done)

- `docker compose up --build` → API healthy, DB ready, migrations applied.
  
- `GET /healthz` → `{ "ok": true }`.
  
- CRUD devices + readings end-to-end.
  
- Alerts computed + serviced successfully.
  
- CI green; container runs as **non-root**.
  
- Frontend web app can show system health + CRUD devices.