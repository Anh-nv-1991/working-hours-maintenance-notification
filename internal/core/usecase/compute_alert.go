package usecase

import (
	"context"
	"maint/internal/domain"
	"maint/internal/ports"
)

type ComputeAlertUseCase struct {
	repo     ports.Repository
	notifier ports.Notifier // optional
}

func NewComputeAlertUseCase(r ports.Repository, n ports.Notifier) *ComputeAlertUseCase {
	return &ComputeAlertUseCase{repo: r, notifier: n}
}

type ComputeAlertInput struct {
	DeviceID int64 `json:"device_id" binding:"required"`
}

type ComputeAlertOutput struct {
	Alert   *domain.Alert `json:"alert,omitempty"`
	Created bool          `json:"created"`
}

// Execute: tính alert từ last reading; tạo alert nếu cần & chưa có alert open.
func (uc *ComputeAlertUseCase) Execute(ctx context.Context, in ComputeAlertInput) (*ComputeAlertOutput, error) {
	out := &ComputeAlertOutput{Created: false}

	// 1) Lấy last reading
	last, err := uc.repo.GetLastReading(ctx, in.DeviceID)
	if err != nil || last == nil {
		// không có reading -> không tạo alert
		return out, err
	}

	// 2) Lấy plan
	plan, err := uc.repo.GetPlanByDevice(ctx, in.DeviceID)
	if err != nil || plan == nil {
		return out, err
	}

	// 3) Check ngưỡng
	msg := plan.OutOfRangeMessage(last.Value)
	if msg == "" {
		return out, nil // trong ngưỡng
	}

	// 4) Nếu đã có alert open -> trả về cái đang open (không tạo mới)
	if open, _ := uc.repo.GetOpenAlertByDevice(ctx, in.DeviceID); open != nil && open.IsOpen() {
		out.Alert = open
		return out, nil
	}

	// 5) Tạo alert mới gắn với last reading
	alert, err := uc.repo.CreateAlert(ctx, in.DeviceID, &last.ID, msg)
	if err != nil {
		return nil, err
	}
	out.Alert = alert
	out.Created = true

	// 6) Notify (best-effort)
	if uc.notifier != nil && out.Alert != nil {
		_ = uc.notifier.NotifyAlert(ctx, out.Alert)
	}

	return out, nil
}
