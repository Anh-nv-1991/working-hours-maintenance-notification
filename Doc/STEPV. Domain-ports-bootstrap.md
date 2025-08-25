# 🚀 CHEATSHEET STEP 5 – DOMAIN / PORTS

🎯 Mục tiêu:  
- Tách domain struct thuần Go khỏi Gin/DB  
- Định nghĩa ports (interfaces) để use-case gọi → dễ test/mock  

---

📂 Cấu trúc thư mục + Shell Commands (PowerShell):
```powershell
cd your_project

mkdir internal\domain
mkdir internal\ports

ni internal\domain\device.go
ni internal\domain\reading.go
ni internal\domain\plan.go
ni internal\domain\alert.go

ni internal\ports\repo.go
ni internal\ports\notifier.go
// internal/domain/device.go
package domain
import "time"
type Device struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// internal/domain/reading.go
package domain
import "time"
type Reading struct {
	ID       int64     `json:"id"`
	DeviceID int64     `json:"device_id"`
	Value    float64   `json:"value"`
	TakenAt  time.Time `json:"taken_at"`
}

// internal/domain/plan.go
package domain
type Plan struct {
	ID          int64   `json:"id"`
	DeviceID    int64   `json:"device_id"`
	Threshold   float64 `json:"threshold"`
	Description string  `json:"description"`
}

// internal/domain/alert.go
package domain
import "time"
type Alert struct {
	ID         int64      `json:"id"`
	DeviceID   int64      `json:"device_id"`
	Message    string     `json:"message"`
	Active     bool       `json:"active"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
// internal/ports/repo.go
package ports
import "module/path/internal/domain"

type DeviceRepo interface {
	Create(device domain.Device) (domain.Device, error)
	GetByID(id int64) (domain.Device, error)
	List() ([]domain.Device, error)
}

type ReadingRepo interface {
	Add(reading domain.Reading) (domain.Reading, error)
	GetLastByDevice(deviceID int64) (domain.Reading, error)
}

type PlanRepo interface {
	GetByDevice(deviceID int64) ([]domain.Plan, error)
}

type AlertRepo interface {
	Create(alert domain.Alert) (domain.Alert, error)
	ListActive(deviceID int64) ([]domain.Alert, error)
	MarkResolved(alertID int64) error
}

// internal/ports/notifier.go
package ports
import "module/path/internal/domain"
type Notifier interface {
	SendAlert(alert domain.Alert) error
}
go fmt ./...
go vet ./...
go build ./...
