package router

import (
	"wh-ma/internal/adapter/inbound/http/handler"

	"github.com/gin-gonic/gin"
)

func MountDevices(rg *gin.RouterGroup, h *handler.DevicesHandler) {
	g := rg.Group("/devices")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.PATCH("/:id", h.UpdateBasic)
	g.PATCH("/:id/plan", h.UpdatePlan)
	g.DELETE("/:id", h.SoftDelete)
}
