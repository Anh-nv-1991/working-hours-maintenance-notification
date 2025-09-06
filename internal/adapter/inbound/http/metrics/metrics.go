package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP
var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Latency distributions of HTTP requests.",
			Buckets: prometheus.DefBuckets, // 10ms..10s
		},
		[]string{"path", "method"},
	)
)

// Domain-specific: devices
var (
	DeviceCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "devices_created_total",
			Help: "Number of devices created.",
		},
	)

	DeviceListTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "devices_list_total",
			Help: "Number of device list operations.",
		},
	)
)
