package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"maint/internal/ports"

	"github.com/gin-gonic/gin"
)

type DevicesHandler struct {
	repo ports.DeviceRepo
}

func NewDevicesHandler(r ports.DeviceRepo) *DevicesHandler {
	return &DevicesHandler{repo: r}
}

func (h *DevicesHandler) PostDevice(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dev, err := h.repo.CreateDevice(c, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "device name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dev)
}

func (h *DevicesHandler) GetDevice(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	dev, err := h.repo.GetDevice(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dev)
}
