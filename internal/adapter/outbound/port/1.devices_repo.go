package port

import (
	"context"
	"time"
	"wh-ma/internal/domain"
)

// DeviceRepository là hợp đồng cho tầng Usecase gọi ra ngoài
type DeviceRepository interface {
	// Tạo mới device
	Create(ctx context.Context, in CreateDeviceInput) (*domain.Device, error)

	// Lấy chi tiết device theo ID
	GetByID(ctx context.Context, id domain.DeviceID) (*domain.Device, error)

	// Danh sách device (có phân trang)
	List(ctx context.Context, limit, offset int32) ([]*domain.Device, error)

	// Update thông tin cơ bản (tên, trạng thái, vị trí)
	UpdateBasic(ctx context.Context, id domain.DeviceID, name string, status domain.DeviceStatus, location *string) (*domain.Device, error)

	// Gán/bỏ Plan cho device
	UpdatePlan(ctx context.Context, id domain.DeviceID, planID *domain.PlanID) (*domain.Device, error)

	// Xóa mềm
	SoftDelete(ctx context.Context, id domain.DeviceID) error
}

// ==== Input struct cho Create ====
// Tách riêng khỏi domain.Device để tránh nhầm lẫn khi map DB <-> Domain
type CreateDeviceInput struct {
	SerialNumber             string
	Name                     string
	Model                    string
	Manufacturer             string
	Year                     int
	CommissionDate           *time.Time // alias cho *time.Time hoặc wrapper domain
	TotalWorkingHour         int
	AfterOverhaulWorkingHour int
	Status                   domain.DeviceStatus
	LastServiceAt            *time.Time
	Location                 string
	PlanID                   *domain.PlanID
}
