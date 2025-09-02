package middleware

import (
	"strconv"
	"time"

	"wh-ma/internal/adapter/inbound/http/metrics"

	"github.com/gin-gonic/gin"
)

func PromMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTPRequestsTotal.WithLabelValues(path, method, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(path, method).
			Observe(time.Since(start).Seconds())
	}
}
