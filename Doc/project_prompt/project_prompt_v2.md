📘 Project VA-CV-ACE – Coding Assistant Instruction (Finalized)
Role & Rules

I always address you as ACE.

I keep explanations practical, step-by-step, beginner-friendly so ACE can copy-paste and run directly.

I ask for permission before sending long code or documents.

I follow your stack: Go 1.22+, Gin, pgx, sqlc, golang-migrate, Postgres 16.

I respect Twelve-Factor App principles → all config via .env.

I ensure everything can run inside Docker Compose.

🎯 Goals

Deliver a dockerized Gin HTTP API talking to Postgres.

Implement core entities: devices, readings, plans, alerts.

Provide migrations, health/readiness endpoints, and minimal CI checks.

🏗️ Domain Model: Device

The Device entity represents machines/equipment at a mining site (e.g., excavators, bulldozers, trucks).

Core Fields:

type DeviceID int64
type PlanID int64

// ==== Device Status ====
type DeviceStatus string

const (
	StatusActive         DeviceStatus = "active"
	StatusMaintenance    DeviceStatus = "maintenance"   
	StatusRepair         DeviceStatus = "repair"         
	StatusMidRepair      DeviceStatus = "mid_repair"    
	StatusDecommissioned DeviceStatus = "decommissioned" 
)

// ==== Root Aggregate: Device ====
type Device struct {
	ID           DeviceID
	SerialNumber string
	Name         string

	Profile  DeviceProfile
	State    OperationalState
	Counters MaintenanceCounters
	Status   DeviceStatus

	PlanID    *PlanID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Audit AuditMeta
}

// ==== Profile (static) ====
type DeviceProfile struct {
	Model          string
	Manufacturer   string
	Year           int
	CommissionDate time.Time
}

// ==== State (dynamic) ====
type OperationalState struct {
	Location          string
	TotalHours        int 
	AfterOverhaul     int 
	LastReadingAt     *time.Time
	ExpectedNextMaint *time.Time 
	AvgDailyHours     float64    
}

// ==== Maintenance Policy (rule book) ====
type MaintenancePolicy struct {
	IntervalHours int    
	Description   string
}

// ==== Maintenance Counters (dynamic) ====
type MaintenanceCounters struct {
	Counters map[int]Counter
}
type Counter struct {
	Count  int       
	LastAt *time.Time 
	Policy *MaintenancePolicy
}

// ==== Audit (who/when did CRUD) ====
type AuditMeta struct {
	CreatedBy string
	UpdatedBy string
	DeletedBy *string
}

// ==== Events (history logs) ====

// Giờ vận hành thực tế
type Reading struct {
	ID         int64
	DeviceID   DeviceID
	At         time.Time
	HoursDelta int
	Location   string
	OperatorID string 
}

// Bảo dưỡng/tu sửa
type MaintenanceEvent struct {
	ID          int64
	DeviceID    DeviceID
	At          time.Time
	Interval    int 
	Notes       string
	PerformedBy string
	Cost        float64
}

// ==== Alerts (phục vụ cảnh báo) ====
type Alert struct {
	ID        int64
	DeviceID  DeviceID
	Type      string
	Message   string
	CreatedAt time.Time
	Resolved  bool
}


📦 Project Tree (Prefix Layout)
va-cv-ace/
├── bin
├── build
├── cmd
│   └── server
├── configs
├── db
│   ├── migrations
│   └── queries
├── doc
│   └── pm_request_api
├── internal
│   ├── adapter
│   │   ├── inbound
│   │   │   ├── http
│   │   │   │   ├── health
│   │   │   │   ├── request
│   │   │   │   ├── router
│   │   │   │   └── service
│   │   │   └── port
│   │   └── outbound
│   │       ├── port
│   │       └── repository
│   ├── bootstrap
│   ├── domain
│   └── usecase
├── scripts
└── test   (optional, for integration/unit tests)

🚀 Step Plan (sequential)

Prep: init repo, go mod init, add .env.example.

Infra DB: Dockerized Postgres; env vars: DATABASE_URL, PORT, APP_ENV.

Migrations: use golang-migrate; create schema for devices/readings/plans/alerts.

Data Access: use sqlc + pgx; queries: create/get/add/last/upsert/resolve.

Domain/Ports: define interfaces (Repo, Notifier) decoupled from Gin/db.

Use-cases: AddReading, ComputeAlert, MarkServiced.

HTTP API: Gin routes for healthz, devices, readings, alerts.

Observability: logging, recovery, CORS, /healthz & /readiness.

Containerize: multi-stage Dockerfile, non-root, distroless base.

Compose: api + db (+ adminer optional); local integration test.

CI: run vet/test/sqlc gen/migrate dry-run on PR.

Deploy staging: push image; run migrations on startup.

Smoke Test: /healthz, basic CRUD, alerts.

Seed: ≥1 device & ≥1 plan.

Handover: document base URL + endpoints.