package httpserver

import (
	"maint/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

type RouterDeps struct {
	Devices         *handlers.DevicesHandler
	Readings        *handlers.ReadingsHandler // nếu có
	Plans           *handlers.PlansHandler
	Alerts          *handlers.AlertsHandler        // nếu dùng biến thể repo-based
	AlertsComputeUC *handlers.AlertsComputeHandler // nếu dùng biến thể usecase-based
}

func NewRouter(d RouterDeps) *gin.Engine {
	r := gin.Default()

	// health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	// devices
	if d.Devices != nil {
		r.POST("/devices", d.Devices.PostDevice)
		r.GET("/devices/:id", d.Devices.GetDevice)
	}

	// readings (bật nếu đã có handler)
	if d.Readings != nil {
		r.POST("/readings", d.Readings.PostReading)
		r.GET("/readings/last/:deviceID", d.Readings.GetLastReading)
	}

	// plans — ĐÃ SỬA TÊN METHOD CHO KHỚP HANDLER
	if d.Plans != nil {
		r.POST("/plans", d.Plans.UpsertPlan)        // ← khớp UpsertPlan
		r.GET("/plans/:device_id", d.Plans.GetPlan) // ← khớp GetPlan
	}

	// alerts (tùy bạn dùng biến thể nào)
	if d.AlertsComputeUC != nil {
		r.POST("/alerts/compute", d.AlertsComputeUC.PostCompute)
	}
	if d.Alerts != nil {
		r.POST("/alerts/compute/:device_id", d.Alerts.ComputeAlert)
		r.POST("/alerts/:id/service", d.Alerts.MarkServiced)
	}

	return r
}
