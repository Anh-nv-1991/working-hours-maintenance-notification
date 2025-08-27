// internal/usecase/add_reading.go
package usecase

import (
	"context"
	"time"

	db "maint/internal/db/sqlc"
)

// 1) Định nghĩa interface tối thiểu mà ComputeAlert cần có
type AlertComputer interface {
	Execute(ctx context.Context, in ComputeAlertInput) (ComputeAlertOutput, error)
}

type AddReadingInput struct {
	DeviceID int64
	Value    float64
	At       *time.Time
}

type AddReadingOutput struct {
	ReadingID int64
}

type AddReadingUseCase struct {
	ReadingsRepo interface {
		AddReading(ctx context.Context, arg db.AddReadingParams) (db.Reading, error)
	}
	// 2) Inject “bộ tính alert” thay vì tự gọi NewComputeAlertUseCase(...)
	AlertComp AlertComputer // có thể nil nếu không muốn auto-compute
}

func NewAddReadingUseCase(
	readingsRepo interface {
		AddReading(ctx context.Context, arg db.AddReadingParams) (db.Reading, error)
	},
	alertComp AlertComputer, // cho phép truyền nil
) *AddReadingUseCase {
	return &AddReadingUseCase{ReadingsRepo: readingsRepo, AlertComp: alertComp}
}

func (uc *AddReadingUseCase) Execute(ctx context.Context, in AddReadingInput) (AddReadingOutput, error) {
	var at any
	if in.At != nil {
		at = *in.At
	} else {
		at = nil // để DB tự NOW()
	}

	// map đúng theo sqlc: AddReadingParams có Column3 là at (or NOW) :contentReference[oaicite:0]{index=0}
	r, err := uc.ReadingsRepo.AddReading(ctx, db.AddReadingParams{
		DeviceID: in.DeviceID,
		Value:    in.Value,
		Column3:  at,
	})
	if err != nil {
		return AddReadingOutput{}, err
	}

	// nếu có inject bộ compute alert thì gọi luôn
	if uc.AlertComp != nil {
		// Fire & forget (bỏ qua error) hoặc bạn có thể bắt lỗi tuỳ yêu cầu
		_, _ = uc.AlertComp.Execute(ctx, ComputeAlertInput{DeviceID: in.DeviceID})
	}

	return AddReadingOutput{ReadingID: r.ID}, nil
}
