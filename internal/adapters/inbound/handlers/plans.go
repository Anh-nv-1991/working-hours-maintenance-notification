package handlers

import (
	"net/http"
	"strconv"
	"time"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"

	"github.com/gin-gonic/gin"
)

type PlansHandler struct {
	repo ports.PlanRepo
}

func NewPlansHandler(r ports.PlanRepo) *PlansHandler {
	return &PlansHandler{repo: r}
}

// POST /plans
func (h *PlansHandler) UpsertPlan(c *gin.Context) {
	var req struct {
		DeviceID     int64      `json:"device_id" binding:"required"`
		ThresholdMin float64    `json:"threshold_min" binding:"required"`
		ThresholdMax float64    `json:"threshold_max" binding:"required"`
		CreatedAt    *time.Time `json:"created_at,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var created interface{}
	if req.CreatedAt != nil {
		created = *req.CreatedAt
	} else {
		created = nil // NOW()
	}

	plan, err := h.repo.UpsertPlan(c, db.UpsertPlanParams{
		DeviceID:     req.DeviceID,
		ThresholdMin: req.ThresholdMin,
		ThresholdMax: req.ThresholdMax,
		Column4:      created, // <— đúng theo sqlc
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}

// GET /plans/:device_id
func (h *PlansHandler) GetPlan(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}
	plan, err := h.repo.GetPlanForDevice(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plan)
}
