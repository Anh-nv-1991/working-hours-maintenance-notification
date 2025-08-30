package repository

import (
	"context"
	"errors"
	"time"
	"wh-ma/internal/adapter/outbound/port"
	dbsqlc "wh-ma/internal/adapter/outbound/repository/sqlc"
	"wh-ma/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ==== Port-facing interface gợi ý (nếu ACE chưa có) ====
// Nên khai báo trong internal/adapter/outbound/port/device_repository.go
// type DeviceRepository interface {
// 	Create(ctx context.Context, in CreateDeviceInput) (*domain.Device, error)
// 	GetByID(ctx context.Context, id domain.DeviceID) (*domain.Device, error)
// 	List(ctx context.Context, limit, offset int32) ([]*domain.Device, error)
// 	UpdateBasic(ctx context.Context, id domain.DeviceID, name string, status domain.DeviceStatus, location *string) (*domain.Device, error)
// 	UpdatePlan(ctx context.Context, id domain.DeviceID, planID *domain.PlanID) (*domain.Device, error)
// 	SoftDelete(ctx context.Context, id domain.DeviceID) error
// }

// ==== Input cho Create (phù hợp Twelve-Factor, tách khỏi domain nếu cần bind JSON) ====

type DeviceRepositoryPG struct {
	pool *pgxpool.Pool
	q    *dbsqlc.Queries
}

func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepositoryPG {
	return &DeviceRepositoryPG{
		pool: pool,
		q:    dbsqlc.New(pool),
	}
}

// ==== Create ====
func (r *DeviceRepositoryPG) Create(ctx context.Context, in port.CreateDeviceInput) (*domain.Device, error) {
	var planID *int64
	if in.PlanID != nil {
		v := int64(*in.PlanID)
		planID = &v
	}
	row, err := r.q.CreateDevice(ctx, dbsqlc.CreateDeviceParams{
		SerialNumber:             in.SerialNumber,
		Name:                     in.Name,
		Model:                    strPtr(in.Model),
		Manufacturer:             strPtr(in.Manufacturer),
		YearOfManufacture:        int32Ptr(in.Year),
		CommissionDate:           dateFromPtr(in.CommissionDate),
		TotalWorkingHour:         int32Ptr(in.TotalWorkingHour),
		AfterOverhaulWorkingHour: int32Ptr(in.AfterOverhaulWorkingHour),
		Status:                   string(in.Status),
		LastServiceAt:            timestamptzFromPtr(in.LastServiceAt),
		Location:                 strPtr(in.Location),
		PlanID:                   planID,
	})
	if err != nil {
		return nil, err
	}
	d := mapSqlcDeviceToDomain(row)
	return &d, nil
}

// ==== Get ====
func (r *DeviceRepositoryPG) GetByID(ctx context.Context, id domain.DeviceID) (*domain.Device, error) {
	row, err := r.q.GetDevice(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	d := mapSqlcDeviceToDomain(row)
	return &d, nil
}

// ==== List (phân trang đơn giản) ====
func (r *DeviceRepositoryPG) List(ctx context.Context, limit, offset int32) ([]*domain.Device, error) {
	rows, err := r.q.ListDevices(ctx, dbsqlc.ListDevicesParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Device, 0, len(rows))
	for _, row := range rows {
		d := mapSqlcDeviceToDomain(row)
		out = append(out, &d)
	}
	return out, nil
}

// ==== UpdateBasic (đổi tên, trạng thái, vị trí) ====
func (r *DeviceRepositoryPG) UpdateBasic(ctx context.Context, id domain.DeviceID, name string, status domain.DeviceStatus, location *string) (*domain.Device, error) {
	row, err := r.q.UpdateDeviceBasic(ctx, dbsqlc.UpdateDeviceBasicParams{
		ID:       int64(id),
		Name:     name,
		Status:   string(status),
		Location: location, // nullable
	})
	if err != nil {
		return nil, err
	}
	d := mapSqlcDeviceToDomain(row)
	return &d, nil
}

// ==== UpdatePlan (gán/bỏ plan) ====
func (r *DeviceRepositoryPG) UpdatePlan(ctx context.Context, id domain.DeviceID, planID *domain.PlanID) (*domain.Device, error) {
	var pid *int64
	if planID != nil {
		v := int64(*planID)
		pid = &v
	}
	row, err := r.q.UpdateDevicePlan(ctx, dbsqlc.UpdateDevicePlanParams{
		ID:     int64(id),
		PlanID: pid,
	})
	if err != nil {
		return nil, err
	}
	d := mapSqlcDeviceToDomain(row)
	return &d, nil
}

// ==== SoftDelete ====
func (r *DeviceRepositoryPG) SoftDelete(ctx context.Context, id domain.DeviceID) error {
	if id == 0 {
		return errors.New("invalid id")
	}
	return r.q.SoftDeleteDevice(ctx, int64(id))
}

// ==== Mapper: sqlc.Device -> domain.Device ====
func mapSqlcDeviceToDomain(x dbsqlc.Device) domain.Device {
	// plan_id -> *domain.PlanID
	var planID *domain.PlanID
	if x.PlanID != nil {
		v := domain.PlanID(*x.PlanID)
		planID = &v
	}

	// pgtype.Date -> time.Time
	var commission time.Time
	if x.CommissionDate.Valid {
		commission = x.CommissionDate.Time
	}

	// pgtype.Timestamptz -> *time.Time
	var lastServiceAt *time.Time
	if x.LastServiceAt.Valid {
		t := x.LastServiceAt.Time
		lastServiceAt = &t
	}
	var expectedNext *time.Time
	if x.ExpectedNextMaint.Valid {
		t := x.ExpectedNextMaint.Time
		expectedNext = &t
	}
	var deletedAt *time.Time
	if x.DeletedAt.Valid {
		t := x.DeletedAt.Time
		deletedAt = &t
	}

	return domain.Device{
		ID:           domain.DeviceID(x.ID),
		SerialNumber: x.SerialNumber,
		Name:         x.Name,

		Profile: domain.DeviceProfile{
			Model:          strOrEmpty(x.Model),
			Manufacturer:   strOrEmpty(x.Manufacturer),
			Year:           int(i32OrZero(x.YearOfManufacture)),
			CommissionDate: commission,
		},
		State: domain.OperationalState{
			Location:          strOrEmpty(x.Location),
			TotalHours:        int(i32OrZero(x.TotalWorkingHour)),
			AfterOverhaul:     int(i32OrZero(x.AfterOverhaulWorkingHour)),
			LastReadingAt:     lastServiceAt,
			ExpectedNextMaint: expectedNext,
			AvgDailyHours:     f64OrZero(x.AvgDailyHours), // <-- FIX: *float64 -> float64
		},
		Status:    domain.DeviceStatus(x.Status),
		PlanID:    planID,
		CreatedAt: x.CreatedAt.Time,
		UpdatedAt: x.UpdatedAt.Time,
		DeletedAt: deletedAt,

		Audit: domain.AuditMeta{
			CreatedBy: strOrEmpty(x.CreatedBy), // <-- FIX: *string -> string
			UpdatedBy: strOrEmpty(x.UpdatedBy),
			DeletedBy: x.DeletedBy, // *string giữ nguyên
		},
	}
}

// ---- helpers an toàn cho con trỏ ----
func strOrEmpty(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}

func f64OrZero(p *float64) float64 {
	if p != nil {
		return *p
	}
	return 0
}

func i32OrZero(p *int32) int32 {
	if p != nil {
		return *p
	}
	return 0
}
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int32Ptr(i int) *int32 {
	if i == 0 {
		return nil
	}
	v := int32(i)
	return &v
}

func dateFromPtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func timestamptzFromPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}
