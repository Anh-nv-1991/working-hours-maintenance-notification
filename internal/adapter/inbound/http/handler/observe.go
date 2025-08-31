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

// getLogger: lấy logger từ middleware (đúng chữ ký) với fallback base logger
func getLogger(ctx context.Context) *slog.Logger {
	base := slog.Default() // hoặc middleware.NewBaseLogger() nếu ACE đã có
	return middleware.GetLogger(ctx, base)
}

func observe(c *gin.Context, op string) (done func(attrs ...slog.Attr)) {
	logger := getLogger(c.Request.Context())
	start := time.Now()

	logger.Info("start "+op,
		slog.String("method", c.Request.Method),
		slog.String("path", c.FullPath()),
	)

	return func(attrs ...slog.Attr) {
		base := []slog.Attr{
			slog.String("op", op),
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Duration("took", time.Since(start)),
		}
		all := append(base, attrs...) // []slog.Attr
		logger.Info("done "+op, toArgs(all)...)
	}
}
