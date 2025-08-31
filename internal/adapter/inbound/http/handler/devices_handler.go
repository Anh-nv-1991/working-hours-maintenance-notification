package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"wh-ma/internal/adapter/inbound/http/metrics"
	"wh-ma/internal/adapter/inbound/http/request"
	inport "wh-ma/internal/adapter/inbound/port"
	"wh-ma/internal/domain"
	"wh-ma/internal/usecase/dto"
)

type DevicesHandler struct {
	svc inport.DevicesInbound
}

func NewDevicesHandler(svc inport.DevicesInbound) *DevicesHandler {
	return &DevicesHandler{svc: svc}
}

// POST /devices
func (h *DevicesHandler) Create(c *gin.Context) {
	done := observe(c, "CreateDevice")
	status := http.StatusCreated
	var errMsg string
	defer func() {
		done(slog.Int("status", status), slog.String("error", errMsg))
	}()
	var in request.CreateDevice
	if err := c.ShouldBindJSON(&in); err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	var planID *domain.PlanID
	if in.PlanID != nil {
		v := domain.PlanID(*in.PlanID)
		planID = &v
	}
	cmd := dto.CreateDeviceCmd{
		SerialNumber:   in.SerialNumber,
		Name:           in.Name,
		Model:          in.Model,
		Manufacturer:   in.Manufacturer,
		Year:           in.Year,
		CommissionDate: in.CommissionDate,
		Status:         domain.DeviceStatus(in.Status),
		Location:       in.Location,
		PlanID:         planID,
	}
	dev, err := h.svc.Create(c, cmd)
	if err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	metrics.DeviceCreatedTotal.Inc()
	c.JSON(http.StatusCreated, dev)
}

// GET /devices/:id
func (h *DevicesHandler) Get(c *gin.Context) {
	done := observe(c, "GetDevice")
	status := http.StatusOK
	var errMsg string
	var id domain.DeviceID

	defer func() {
		done(
			slog.Int("status", status),
			slog.String("error", errMsg),
			slog.Int64("device_id", int64(id)),
		)
	}()

	var ok bool
	id, ok = parseDeviceID(c)
	if !ok {
		status = http.StatusBadRequest
		errMsg = "invalid id"
		return // parseDeviceID đã trả JSON 400 rồi
	}

	dev, err := h.svc.Get(c, id)
	if err != nil {
		status = http.StatusNotFound
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(status, dev)
}

// GET /devices
func (h *DevicesHandler) List(c *gin.Context) {
	done := observe(c, "ListDevices")
	status := http.StatusOK
	var errMsg string
	limit, offset := parsePaging(c, 50, 0)

	defer func() {
		done(
			slog.Int("status", status),
			slog.String("error", errMsg),
			slog.Int("limit", int(limit)),
			slog.Int("offset", int(offset)),
		)
	}()

	devs, err := h.svc.List(c, limit, offset)
	if err != nil {
		status = http.StatusInternalServerError
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	metrics.DeviceListTotal.Inc()
	c.JSON(status, gin.H{"items": devs, "limit": limit, "offset": offset})
}

// PATCH /devices/:id
func (h *DevicesHandler) UpdateBasic(c *gin.Context) {
	done := observe(c, "UpdateDeviceBasic")
	status := http.StatusOK
	var errMsg string
	var id domain.DeviceID

	defer func() {
		done(
			slog.Int("status", status),
			slog.String("error", errMsg),
			slog.Int64("device_id", int64(id)),
		)
	}()

	var ok bool
	id, ok = parseDeviceID(c)
	if !ok {
		status = http.StatusBadRequest
		errMsg = "invalid id"
		return
	}

	var in request.UpdateBasic
	if err := c.ShouldBindJSON(&in); err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	cmd := dto.UpdateDeviceBasicCmd{
		ID:       id,
		Name:     in.Name,
		Status:   domain.DeviceStatus(in.Status),
		Location: in.Location,
	}
	dev, err := h.svc.UpdateBasic(c, cmd)
	if err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(status, dev)
}

// PATCH /devices/:id/plan
func (h *DevicesHandler) UpdatePlan(c *gin.Context) {
	done := observe(c, "UpdateDevicePlan")
	status := http.StatusOK
	var errMsg string
	var id domain.DeviceID
	var planID *domain.PlanID

	defer func() {
		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("error", errMsg),
			slog.Int64("device_id", int64(id)),
		}
		if planID != nil {
			attrs = append(attrs, slog.Int64("plan_id", int64(*planID)))
		}
		done(attrs...)
	}()

	var ok bool
	id, ok = parseDeviceID(c)
	if !ok {
		status = http.StatusBadRequest
		errMsg = "invalid id"
		return
	}

	var in request.UpdatePlan
	if err := c.ShouldBindJSON(&in); err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	if in.PlanID != nil {
		v := domain.PlanID(*in.PlanID)
		planID = &v
	}

	cmd := dto.UpdateDevicePlanCmd{ID: id, PlanID: planID}
	dev, err := h.svc.UpdatePlan(c, cmd)
	if err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.JSON(status, dev)
}

// DELETE /devices/:id (soft delete)
func (h *DevicesHandler) SoftDelete(c *gin.Context) {
	done := observe(c, "SoftDeleteDevice")
	status := http.StatusNoContent
	var errMsg string
	var id domain.DeviceID

	defer func() {
		done(
			slog.Int("status", status),
			slog.String("error", errMsg),
			slog.Int64("device_id", int64(id)),
		)
	}()

	var ok bool
	id, ok = parseDeviceID(c)
	if !ok {
		status = http.StatusBadRequest
		errMsg = "invalid id"
		return
	}

	if err := h.svc.SoftDelete(c, id); err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
		c.JSON(status, gin.H{"error": errMsg})
		return
	}
	c.Status(status) // 204
}

// ===== helpers =====
func parseDeviceID(c *gin.Context) (domain.DeviceID, bool) {
	var uri struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return domain.DeviceID(uri.ID), true
}
func parsePaging(c *gin.Context, defLimit, defOffset int32) (int32, int32) {
	var in struct {
		Limit  *int32 `form:"limit"`
		Offset *int32 `form:"offset"`
	}
	_ = c.ShouldBindQuery(&in)

	limit, offset := defLimit, defOffset
	if in.Limit != nil {
		limit = *in.Limit
	}
	if in.Offset != nil {
		offset = *in.Offset
	}

	// (tuỳ chọn) ràng buộc an toàn
	if limit < 1 {
		limit = 1
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
