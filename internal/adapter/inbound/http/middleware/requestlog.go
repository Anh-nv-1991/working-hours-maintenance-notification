package middleware

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type ctxLoggerKey struct{}

const ginCtxLoggerKey = "logger"

// Ghi log cho mỗi request; không crash nếu base == nil.
// Tự đính kèm trace_id/span_id nếu OTel đã init.
// Đồng thời gắn logger vào cả gin.Context lẫn request.Context().
func RequestLogMiddleware(base *slog.Logger) gin.HandlerFunc {
	if base == nil {
		base = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	return func(c *gin.Context) {
		start := time.Now()
		reqID := requestid.Get(c)

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		l := base.With(
			"request_id", reqID,
			"method", c.Request.Method,
			"route", route,
			"ip", c.ClientIP(),
			"ua", c.Request.UserAgent(),
		)

		// Correlate với OTel nếu có span hiện hành
		if span := oteltrace.SpanFromContext(c.Request.Context()); span != nil {
			sc := span.SpanContext()
			if sc.HasTraceID() {
				l = l.With(
					"trace_id", sc.TraceID().String(),
					"span_id", sc.SpanID().String(),
				)
			}
		}

		// Cho phép handler khác lấy logger
		c.Set(ginCtxLoggerKey, l)
		ctx := context.WithValue(c.Request.Context(), ctxLoggerKey{}, l)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)
		l.Info("request.done",
			slog.Int("status", c.Writer.Status()),
			slog.Int("size", c.Writer.Size()),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("err", strings.Join(c.Errors.Errors(), "; ")),
		)
	}
}

// GetLogger trả logger từ *gin.Context (fallback stdout nếu thiếu).
func GetLogger(c *gin.Context) *slog.Logger {
	if v, ok := c.Get(ginCtxLoggerKey); ok {
		if l, ok2 := v.(*slog.Logger); ok2 && l != nil {
			return l
		}
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
}

// CtxLogger trả logger từ context.Context (dùng cho các chỗ chỉ có ctx).
// Nếu không tìm thấy thì trả fallback; nếu fallback nil thì trả stdout logger.
func CtxLogger(ctx context.Context, fallback *slog.Logger) *slog.Logger {
	if ctx != nil {
		if v := ctx.Value(ctxLoggerKey{}); v != nil {
			if l, ok := v.(*slog.Logger); ok && l != nil {
				return l
			}
		}
	}
	if fallback != nil {
		return fallback
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
}
