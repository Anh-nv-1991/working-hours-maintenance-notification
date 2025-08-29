package devicesdomain

import "time"

// ==== Strongly typed IDs ====
type DeviceID int64
type PlanID int64

// ==== Device Status ====
type DeviceStatus string

const (
	StatusActive         DeviceStatus = "active"
	StatusMaintenance    DeviceStatus = "maintenance"    // bảo dưỡng định kỳ
	StatusRepair         DeviceStatus = "repair"         // sửa chữa đột xuất
	StatusMidRepair      DeviceStatus = "mid_repair"     // trung tu
	StatusDecommissioned DeviceStatus = "decommissioned" // ngừng/loại bỏ
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
	TotalHours        int // lifetime TWH
	AfterOverhaul     int // AOHWH
	LastReadingAt     *time.Time
	ExpectedNextMaint *time.Time // dự đoán lần bảo dưỡng kế tiếp
	AvgDailyHours     float64    // trung bình giờ/ngày
}

// ==== Maintenance Policy (rule book) ====
type MaintenancePolicy struct {
	IntervalHours int    // e.g. 250, 500, 1000
	Description   string // thay dầu, kiểm tra phanh...
}

// ==== Maintenance Counters (dynamic) ====
type MaintenanceCounters struct {
	// map interval_hours -> Counter
	Counters map[int]Counter
}

type Counter struct {
	Count  int        // số lần đã làm
	LastAt *time.Time // thời gian gần nhất
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
	OperatorID string // bổ sung để phân tích hành vi người vận hành
}

// Bảo dưỡng/tu sửa
type MaintenanceEvent struct {
	ID          int64
	DeviceID    DeviceID
	At          time.Time
	Interval    int // ví dụ 250, 500, 1000...
	Notes       string
	PerformedBy string
	Cost        float64
}

// ==== Alerts (phục vụ cảnh báo) ====
type Alert struct {
	ID        int64
	DeviceID  DeviceID
	Type      string // "maintenance_due", "over_usage", "idle_too_long"
	Message   string
	CreatedAt time.Time
	Resolved  bool
}
