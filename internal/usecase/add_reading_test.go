package usecase_test

import (
	"context"
	"testing"

	"working_hours/internal/domain"
)

type fakeClock struct{ t int64 }

func (f fakeClock) NowUnix() int64 { return f.t }

type fakeRepo struct {
	withTx func(ctx context.Context, fn func(r usecase_repo) error) error
	// implement các method tối thiểu; hoặc dùng gomock/testify nếu ACE thích
}

// … (để ngắn gọn, em sẽ viết full test nếu ACE cần)

func TestAddReading_CreateAlertWhenOutOfRange(t *testing.T) {
	// Ý tưởng: fake repo trả plan {min:10,max:50}, add reading 60 => tạo alert
	// Kiểm tra Execute trả out.Alert != nil
	// Em có thể bổ sung full mock nếu ACE muốn.
	_ = domain.Device{} // tránh unused; placeholder
}
