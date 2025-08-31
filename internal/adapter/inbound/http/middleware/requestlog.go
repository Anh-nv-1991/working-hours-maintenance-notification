package middleware

import (
	"context"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const LoggerKey ctxKey = "logger"

// NewBaseLogger tạo logger gốc theo LOG_LEVEL (env) và output JSON.
func NewBaseLogger() *slog.Logger {
	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

// RequestLogMiddleware:
// - đảm bảo có X-Request-ID (tạo nếu client không gửi)
// - đưa logger vào context, đính kèm field cơ bản
// - log access line ở cuối (status, latency, size)
func RequestLogMiddleware(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// đảm bảo đã có request id
		reqID := requestid.Get(c)
		start := time.Now()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		ip := clientIP(c)
		userAgent := c.Request.UserAgent()

		// logger cho request hiện tại
		l := base.With(
			"request_id", reqID,
			"method", c.Request.Method,
			"route", route,
			"ip", ip,
		)

		// đưa logger vào context chuẩn
		ctx := context.WithValue(c.Request.Context(), LoggerKey, l)
		c.Request = c.Request.WithContext(ctx)

		// log start ở mức debug (không ồn khi level >= info)
		l.Debug("request.start", "ua", userAgent)

		// chạy các handler tiếp theo
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// access log: 1 dòng/tác vụ
		l.Info("request.done",
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"size", size,
			"err", strings.Join(c.Errors.Errors(), "; "),
		)
	}
}

// GetLogger rút logger từ context. Nếu thiếu, trả base để không nil.
func GetLogger(ctx context.Context, fallback *slog.Logger) *slog.Logger {
	if v := ctx.Value(LoggerKey); v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return fallback
}

// clientIP lấy IP thực tế, tôn trọng X-Forwarded-For nếu có reverse proxy.
func clientIP(c *gin.Context) string {
	// Ưu tiên X-Forwarded-For (chuỗi ip, lấy ip đầu)
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// Fallback: RemoteAddr
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.ClientIP() // gin đã cố gắng parse giúp
	}
	return host
}
