package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
	var in request.CreateDevice
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dev)
}

// GET /devices/:id
func (h *DevicesHandler) Get(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	dev, err := h.svc.Get(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dev)
}

// GET /devices
func (h *DevicesHandler) List(c *gin.Context) {
	limit, offset := parsePaging(c, 50, 0)
	devs, err := h.svc.List(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": devs, "limit": limit, "offset": offset})
}

// PATCH /devices/:id
func (h *DevicesHandler) UpdateBasic(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	var in request.UpdateBasic
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dev)
}

// PATCH /devices/:id/plan
func (h *DevicesHandler) UpdatePlan(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	var in request.UpdatePlan
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var planID *domain.PlanID
	if in.PlanID != nil {
		v := domain.PlanID(*in.PlanID)
		planID = &v
	}
	cmd := dto.UpdateDevicePlanCmd{ID: id, PlanID: planID}
	dev, err := h.svc.UpdatePlan(c, cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dev)
}

// DELETE /devices/:id (soft delete)
func (h *DevicesHandler) SoftDelete(c *gin.Context) {
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	if err := h.svc.SoftDelete(c, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
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
	type q struct {
		Limit, Offset *int32 `form:"limit,offset"`
	}
	var qq struct {
		Limit  *int32 `form:"limit"`
		Offset *int32 `form:"offset"`
	}
	_ = c.ShouldBindQuery(&qq)
	limit, offset := defLimit, defOffset
	if qq.Limit != nil {
		limit = *qq.Limit
	}
	if qq.Offset != nil {
		offset = *qq.Offset
	}
	return limit, offset
}
