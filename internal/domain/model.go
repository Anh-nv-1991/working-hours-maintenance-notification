package domain

import "time"

// Device: thiết bị (máy, cảm biến, ...).
type Device struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Reading: số liệu đo được từ thiết bị.
type Reading struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"device_id"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// Plan: cấu hình ngưỡng cho thiết bị.
type Plan struct {
	ID           int64     `json:"id"`
	DeviceID     int64     `json:"device_id"`
	ThresholdMin float64   `json:"threshold_min"`
	ThresholdMax float64   `json:"threshold_max"`
	CreatedAt    time.Time `json:"created_at"`
}

// Alert: cảnh báo khi reading vượt ngưỡng.
type Alert struct {
	ID         int64      `json:"id"`
	DeviceID   int64      `json:"device_id"`
	ReadingID  *int64     `json:"reading_id,omitempty"`
	Message    string     `json:"message"`
	Status     string     `json:"status"` // "open" | "serviced"
	CreatedAt  time.Time  `json:"created_at"`
	ServicedAt *time.Time `json:"serviced_at,omitempty"`
}
