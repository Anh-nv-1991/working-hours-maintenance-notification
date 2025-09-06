package handler

import (
	"context"
	"log/slog"
	"time"

	"wh-ma/internal/adapter/inbound/http/middleware"

	"github.com/gin-gonic/gin"
)

// toArgs chuyển []slog.Attr -> []any để truyền vào logger.Info(...any)
func toArgs(attrs []slog.Attr) []any {
	out := make([]any, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, a)
	}
	return out
}

// Lấy logger từ context (đã được middleware gắn); fallback = slog.Default()
func getLogger(ctx context.Context) *slog.Logger {
	return middleware.CtxLogger(ctx, slog.Default())
}

func observe(c *gin.Context, op string) (done func(attrs ...slog.Attr)) {
	logger := getLogger(c.Request.Context())
	start := time.Now()

	route := c.FullPath()
	if route == "" {
		route = c.Request.URL.Path
	}

	logger.Info("start "+op,
		slog.String("method", c.Request.Method),
		slog.String("path", route),
	)

	return func(attrs ...slog.Attr) {
		base := []slog.Attr{
			slog.String("op", op),
			slog.String("method", c.Request.Method),
			slog.String("path", route),
			slog.Duration("took", time.Since(start)),
		}
		all := append(base, attrs...)
		logger.Info("done "+op, toArgs(all)...)
	}
}
