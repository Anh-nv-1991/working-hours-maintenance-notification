package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	startedAt time.Time
	pool      *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		startedAt: time.Now(),
		pool:      pool,
	}
}

// GET /healthz  -> app liveness (không ping DB để cực nhanh)
func (h *Handler) Liveness(c *gin.Context) {
	uptime := time.Since(h.startedAt).Round(time.Second).String()
	c.JSON(http.StatusOK, gin.H{
		"ok":        true,
		"service":   "wh-ma",
		"uptime":    uptime,
		"timestamp": time.Now().UTC(),
	})
}

// GET /readiness -> readiness (ping DB với timeout ngắn)
func (h *Handler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 1*time.Second)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ok":      false,
			"reason":  "db_unreachable",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"db":      "up",
		"message": "ready",
	})
}
