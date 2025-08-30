package usecase

import (
	"context"
	"errors"
	"time"

	inport "wh-ma/internal/adapter/inbound/port"
	outport "wh-ma/internal/adapter/outbound/port"
	"wh-ma/internal/domain"
	"wh-ma/internal/usecase/dto"
)

type DevicesUsecase struct {
	devRepo   outport.DeviceRepository
	planRepo  outport.PlanRepository
	alertRepo outport.AlertRepository
}

func NewDevicesUsecase(
	devRepo outport.DeviceRepository,
	planRepo outport.PlanRepository,
	alertRepo outport.AlertRepository,
) *DevicesUsecase {
	return &DevicesUsecase{devRepo: devRepo, planRepo: planRepo, alertRepo: alertRepo}
}

// ✅ compile-time check: UC triển khai inbound port
var _ inport.DevicesInbound = (*DevicesUsecase)(nil)

// 1) CREATE
// - required: SerialNumber, Name
// - valid: Year >= 1970, CommissionDate <= today
// - default Status=active if empty
// - if PlanID != nil -> verify plan exists
// - ExpectedNextMaint: chưa tính ở đây (để Readings/Maintenance)
func (uc *DevicesUsecase) Create(ctx context.Context, in dto.CreateDeviceCmd) (*domain.Device, error) {
	if in.SerialNumber == "" {
		return nil, errors.New("serial_number is required")
	}
	if in.Name == "" {
		return nil, errors.New("name is required")
	}
	if in.Year < 1970 {
		return nil, errors.New("year is invalid (must be >= 1970)")
	}
	if in.CommissionDate != nil && in.CommissionDate.After(time.Now()) {
		return nil, errors.New("commission_date cannot be in the future")
	}
	if in.Status == "" {
		in.Status = domain.StatusActive
	}
	if in.PlanID != nil {
		if _, err := uc.planRepo.GetByID(ctx, *in.PlanID); err != nil {
			return nil, err
		}
	}

	// map DTO -> input repo (outbound)
	repoIn := outport.CreateDeviceInput{
		SerialNumber:   in.SerialNumber,
		Name:           in.Name,
		Model:          in.Model,
		Manufacturer:   in.Manufacturer,
		Year:           in.Year,
		CommissionDate: in.CommissionDate,
		Status:         in.Status,
		Location:       valOrEmpty(in.Location),
		PlanID:         in.PlanID,
	}
	return uc.devRepo.Create(ctx, repoIn)
}

// 5) GET/LIST: thuần repo
func (uc *DevicesUsecase) Get(ctx context.Context, id domain.DeviceID) (*domain.Device, error) {
	return uc.devRepo.GetByID(ctx, id)
}
func (uc *DevicesUsecase) List(ctx context.Context, limit, offset int32) ([]*domain.Device, error) {
	return uc.devRepo.List(ctx, limit, offset)
}

// 2) UPDATE BASIC
// - status chỉ cho phép: active/maintenance/repair/mid_repair/decommissioned
// - CHO phép chuyển từ decommissioned -> active (theo yêu cầu ACE)
// - location có thể nil/"" đều được
func (uc *DevicesUsecase) UpdateBasic(ctx context.Context, in dto.UpdateDeviceBasicCmd) (*domain.Device, error) {
	if in.Name == "" {
		return nil, errors.New("name is required")
	}
	if !isAllowedStatus(in.Status) {
		return nil, errors.New("invalid status")
	}
	return uc.devRepo.UpdateBasic(ctx, in.ID, in.Name, in.Status, in.Location)
}

// 3) UPDATE PLAN
// - gắn plan: verify tồn tại
// - nếu AfterOverhaul >= IntervalHours ở thời điểm gắn -> tạo alert "maintenance_due" nếu chưa có
// - bỏ plan: chỉ ghi nhận, không tạo/đóng alert
func (uc *DevicesUsecase) UpdatePlan(ctx context.Context, in dto.UpdateDevicePlanCmd) (*domain.Device, error) {
	if in.PlanID != nil {
		if _, err := uc.planRepo.GetByID(ctx, *in.PlanID); err != nil {
			return nil, err
		}
	}
	dev, err := uc.devRepo.UpdatePlan(ctx, in.ID, in.PlanID)
	if err != nil {
		return nil, err
	}

	if in.PlanID != nil { // vừa gắn plan
		plan, perr := uc.planRepo.GetByID(ctx, *in.PlanID)
		if perr == nil && plan.IntervalHours > 0 {
			if dev.State.AfterOverhaul >= plan.IntervalHours {
				open, _ := uc.alertRepo.ListOpenByDevice(ctx, dev.ID, 50, 0)
				for _, a := range open {
					if a.Type == "maintenance_due" {
						return dev, nil // đã có alert mở, thôi
					}
				}
				_, _ = uc.alertRepo.Create(ctx, outport.CreateAlertInput{
					DeviceID: dev.ID,
					Type:     "maintenance_due",
					Message:  "Thiết bị đã vượt ngưỡng giờ bảo dưỡng theo kế hoạch mới",
				})
			}
		}
	}
	return dev, nil
}

// 4. SOFT DELETE
//   - chỉ xóa mềm khi KHÔNG còn alert mở
//   - KHÔNG xóa mềm nếu status là maintenance hoặc repair (đang thao tác kỹ thuật)
//     (mid_repair được phép xóa theo yêu cầu)
func (uc *DevicesUsecase) SoftDelete(ctx context.Context, id domain.DeviceID) error {
	dev, err := uc.devRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if dev.Status == domain.StatusMaintenance || dev.Status == domain.StatusRepair {
		return errors.New("cannot soft delete while device is under maintenance/repair")
	}
	open, _ := uc.alertRepo.ListOpenByDevice(ctx, id, 1, 0)
	if len(open) > 0 {
		return errors.New("cannot soft delete while there are open alerts")
	}
	return uc.devRepo.SoftDelete(ctx, id)
}

// --- helpers ---
func isAllowedStatus(s domain.DeviceStatus) bool {
	switch s {
	case domain.StatusActive,
		domain.StatusMaintenance,
		domain.StatusRepair,
		domain.StatusMidRepair,
		domain.StatusDecommissioned:
		return true
	default:
		return false
	}
}
func valOrEmpty(p *string) string {
	if p != nil {
		return *p
	}
	return ""
}
