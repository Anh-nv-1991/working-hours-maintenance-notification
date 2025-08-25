package handlers

import (
	"net/http"

	"maint/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AlertsComputeHandler struct {
	uc *usecase.ComputeAlertUseCase
}

func NewAlertsComputeHandler(uc *usecase.ComputeAlertUseCase) *AlertsComputeHandler {
	return &AlertsComputeHandler{uc: uc}
}

// POST /alerts/compute  body: { "device_id": 1 }
func (h *AlertsComputeHandler) PostCompute(c *gin.Context) {
	var in usecase.ComputeAlertInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	out, err := h.uc.Execute(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	status := http.StatusOK
	if out.Created {
		status = http.StatusCreated
	}
	c.JSON(status, out)
}
