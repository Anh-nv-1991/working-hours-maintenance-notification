package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"wh-ma/internal/adapter/inbound/http/health"
)

func New(p *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	// middlewares: recovery, CORS… (thêm sau nếu cần)

	// Health routes
	h := health.NewHandler(p)
	r.GET("/healthz", h.Liveness)
	r.GET("/readiness", h.Readiness)

	// TODO: các routes khác (devices, readings, alerts, ...)
	return r
}
