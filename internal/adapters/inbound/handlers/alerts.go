package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type AlertsHandler struct {
	repo ports.AlertRepo
}

func NewAlertsHandler(r ports.AlertRepo) *AlertsHandler {
	return &AlertsHandler{repo: r}
}

// POST /alerts/compute/:device_id
func (h *AlertsHandler) ComputeAlert(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	// 1) Check breach
	res, err := h.repo.CheckThresholdBreach(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.Status == "OK" {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
		return
	}

	// 2) Create alert
	alert, err := h.repo.CreateAlert(c, db.CreateAlertParams{
		DeviceID:  id,
		ReadingID: pgtype.Int8{Int64: res.ReadingID, Valid: true},
		Level:     res.Status,
		Message:   pgtype.Text{String: fmt.Sprintf("Value %.2f breached", res.ReadingValue), Valid: true},
		Column5:   nil, // NOW()
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, alert)
}

// POST /alerts/:id/service
func (h *AlertsHandler) MarkServiced(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Cách 1 (đơn giản): để DB tự set NOW() qua COALESCE($2, NOW())
	alert, err := h.repo.ResolveAlert(c, db.ResolveAlertParams{
		ID:         id,
		ServicedAt: pgtype.Timestamptz{}, // Valid=false => NULL => NOW()
	})

	// // Cách 2 (nếu muốn set thời điểm cụ thể):
	// alert, err := h.repo.ResolveAlert(c, db.ResolveAlertParams{
	// 	ID:         id,
	// 	ServicedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	// })

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alert)
}
