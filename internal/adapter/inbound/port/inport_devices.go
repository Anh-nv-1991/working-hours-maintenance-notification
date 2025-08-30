package port

import (
	"context"
	"wh-ma/internal/domain"
	"wh-ma/internal/usecase/dto"
)

type DevicesInbound interface {
	// 1) Create
	Create(ctx context.Context, in dto.CreateDeviceCmd) (*domain.Device, error)

	// 5) Get/List
	Get(ctx context.Context, id domain.DeviceID) (*domain.Device, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.Device, error)

	// 2) UpdateBasic
	UpdateBasic(ctx context.Context, in dto.UpdateDeviceBasicCmd) (*domain.Device, error)

	// 3) UpdatePlan
	UpdatePlan(ctx context.Context, in dto.UpdateDevicePlanCmd) (*domain.Device, error)

	// 4) SoftDelete
	SoftDelete(ctx context.Context, id domain.DeviceID) error
}
