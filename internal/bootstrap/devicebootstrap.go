package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"wh-ma/internal/adapter/inbound/http/handler"
	health "wh-ma/internal/adapter/inbound/http/health"
	"wh-ma/internal/adapter/inbound/http/router"
	outrepo "wh-ma/internal/adapter/outbound/repository"
	"wh-ma/internal/usecase"
)

type AppConfig struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	AllowOrigin []string
}

// LoadConfig: đọc .env qua Compose; không dùng thư viện ngoài để giữ gọn
func LoadConfig() AppConfig {
	cfg := AppConfig{
		AppEnv:      getenv("APP_ENV", "development"),
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", ""),
	}
	origins := getenv("CORS_ORIGINS", "*")
	if origins == "" {
		cfg.AllowOrigin = []string{"*"}
	} else {
		cfg.AllowOrigin = strings.Split(origins, ",")
	}
	return cfg
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// NewPGXPool: khởi tạo pool và Ping để chắc chắn DB sống
func NewPGXPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// BuildRouter: tạo Gin router, middleware, health, mount modules
func BuildRouter(cfg AppConfig, pool *pgxpool.Pool) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS
	c := cors.DefaultConfig()
	if len(cfg.AllowOrigin) == 1 && cfg.AllowOrigin[0] == "*" {
		c.AllowAllOrigins = true
	} else {
		c.AllowOrigins = cfg.AllowOrigin
	}
	c.AllowHeaders = []string{"Authorization", "Content-Type"}
	c.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(c))

	// Health/Readiness
	hh := health.NewHandler(pool)
	r.GET("/healthz", hh.Liveness)    // liveness: không ping DB, trả uptime & timestamp
	r.GET("/readiness", hh.Readiness) // readiness: ping DB với timeout ngắn

	// Outbound repos (PG/sqlc)
	devRepo := outrepo.NewDeviceRepository(pool)
	planRepo := outrepo.NewPlanRepository(pool)
	alertRepo := outrepo.NewAlertRepository(pool)
	// (sẵn sàng cho phần khác) readRepo := outrepo.NewReadingRepository(pool)
	// (sẵn sàng cho phần khác) mntRepo  := outrepo.NewMaintenanceRepository(pool)

	// Usecases (inbound implementations)
	devUC := usecase.NewDevicesUsecase(devRepo, planRepo, alertRepo)

	// Handlers
	devH := handler.NewDevicesHandler(devUC)

	// Router groups
	api := r.Group("/api")
	router.MountDevices(api, devH)

	// Optional: versioning
	// v1 := r.Group("/api/v1")
	// router.MountDevices(v1, devH)

	return r
}

// RunHTTP: chạy server kèm graceful shutdown
func RunHTTP(r *gin.Engine, port string) error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("HTTP listening on :%s", port)
	return srv.ListenAndServe()
}
