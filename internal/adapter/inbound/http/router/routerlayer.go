package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"wh-ma/internal/adapter/inbound/http/health"
	"wh-ma/internal/adapter/inbound/http/middleware"
)

// Options cho Router.New để cấu hình CORS/mode
type Options struct {
	AppEnv      string   // "production" => gin.ReleaseMode
	AllowOrigin []string // []{"*"} => AllowAllOrigins
}

// New tạo *gin.Engine với middleware & infra endpoints
// - Recovery, RequestID, Logger
// - CORS
// - Prometheus middleware + /metrics
// - /healthz, /readiness
//
// Domain endpoints (devices, plans, alerts, readings) sẽ được mount từ bootstrap
// qua các hàm router.Mount* vào group /api.
func New(p *pgxpool.Pool, baseLogger *slog.Logger, opt Options) *gin.Engine {
	if opt.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Middlewares nền tảng
	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(middleware.OTelMiddleware("wh-ma-api"))
	r.Use(middleware.RequestLogMiddleware(baseLogger)) // logger có request-id
	r.Use(middleware.PromMetrics())
	r.Use(middleware.PrometheusHTTP()) // đo count/latency/status cho mọi request

	// CORS
	c := cors.DefaultConfig()
	if len(opt.AllowOrigin) == 1 && opt.AllowOrigin[0] == "*" {
		c.AllowAllOrigins = true
	} else if len(opt.AllowOrigin) > 0 {
		c.AllowOrigins = opt.AllowOrigin
	} else {
		c.AllowAllOrigins = true
	}
	c.AllowHeaders = []string{"Authorization", "Content-Type"}
	c.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(c))

	// Infra endpoints
	h := health.NewHandler(p)
	r.GET("/healthz", h.Liveness)                    // liveness: không ping DB
	r.GET("/readiness", h.Readiness)                 // readiness: ping DB ngắn
	r.GET("/metrics", gin.WrapH(promhttp.Handler())) // Prometheus scrape

	// Ping root (optional)
	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// Domain routes sẽ được mount ở bootstrap:
	//   api := r.Group("/api")
	//   router.MountDevices(api, devicesHandler)
	//   router.MountPlans(api, plansHandler)
	//   router.MountAlerts(api, alertsHandler)
	//   router.MountReadings(api, readingsHandler)

	return r
}
