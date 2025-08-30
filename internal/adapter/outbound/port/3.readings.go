package port

import (
	"context"
	"time"

	"wh-ma/internal/domain"
)

// Hợp đồng để Usecase gọi
type ReadingRepository interface {
	Create(ctx context.Context, in CreateReadingInput) (*domain.Reading, error)
	GetLastByDevice(ctx context.Context, deviceID domain.DeviceID) (*domain.Reading, error)
	ListByDevice(ctx context.Context, deviceID domain.DeviceID, limit, offset int32) ([]*domain.Reading, error)
	Delete(ctx context.Context, id int64) error
}

// Payload tạo mới Reading
type CreateReadingInput struct {
	DeviceID   domain.DeviceID
	At         time.Time
	HoursDelta int
	Location   *string
	OperatorID *string
}
