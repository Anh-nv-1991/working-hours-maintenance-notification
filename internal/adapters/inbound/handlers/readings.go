package handlers

import (
	"net/http"
	"strconv"
	"time"

	db "maint/internal/db/sqlc"
	"maint/internal/ports"

	"github.com/gin-gonic/gin"
)

type ReadingsHandler struct {
	repo ports.ReadingRepo
}

func NewReadingsHandler(r ports.ReadingRepo) *ReadingsHandler {
	return &ReadingsHandler{repo: r}
}

// POST /readings
func (h *ReadingsHandler) PostReading(c *gin.Context) {
	var req struct {
		DeviceID int64      `json:"device_id" binding:"required"`
		Value    float64    `json:"value" binding:"required"`
		At       *time.Time `json:"at,omitempty"` // optional
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var at interface{}
	if req.At != nil {
		at = *req.At
	} else {
		at = nil // NOW()
	}

	reading, err := h.repo.AddReading(c, db.AddReadingParams{
		DeviceID: req.DeviceID,
		Value:    req.Value,
		Column3:  at, // <— đúng theo sqlc
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, reading)
}

// GET /readings/last/:device_id
func (h *ReadingsHandler) GetLastReading(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}
	reading, err := h.repo.LastReading(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reading)
}
