package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func OTelMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		// tên tạm trước khi Gin match route
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		ctx, span := tracer.Start(c.Request.Context(), route,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.route", route),
				attribute.String("http.target", c.Request.URL.RequestURI()),
			),
		)

		// >>> đặt Trace-Id vào response càng sớm càng tốt
		if tid := span.SpanContext().TraceID().String(); tid != "00000000000000000000000000000000" {
			c.Writer.Header().Set("Trace-Id", tid)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// sau khi match xong route, cập nhật lại tên/thuộc tính cho chuẩn
		if rp := c.FullPath(); rp != "" && rp != route {
			span.SetName(rp)
			span.SetAttributes(attribute.String("http.route", rp))
		}

		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
		if len(c.Errors) > 0 || c.Writer.Status() >= 500 {
			span.SetStatus(codes.Error, c.Errors.String())
		} else {
			span.SetStatus(codes.Ok, "OK")
		}
		span.End()
	}
}
