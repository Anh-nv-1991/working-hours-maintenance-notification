package bootstrap

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"wh-ma/internal/adapter/inbound/http/handler"
	"wh-ma/internal/adapter/inbound/http/router"
	outrepo "wh-ma/internal/adapter/outbound/repository"
	"wh-ma/internal/usecase"
)

// ===== Config =====

type AppConfig struct {
	AppEnv      string
	Port        string
	DatabaseURL string
	AllowOrigin []string
	LogLevel    string
}

func LoadConfig() AppConfig {
	LoadEnvFirst()
	cfg := AppConfig{
		AppEnv:      getenv("APP_ENV", "development"),
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", ""),
		LogLevel:    getenv("LOG_LEVEL", "info"),
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

// ===== DB Pool =====

func NewPGXPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

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

// ===== HTTP wiring (router layer định nghĩa endpoints) =====

func BuildRouter(cfg AppConfig, pool *pgxpool.Pool, baseLogger *slog.Logger) *gin.Engine {
	// 1) Repos
	devRepo := outrepo.NewDeviceRepository(pool)
	planRepo := outrepo.NewPlanRepository(pool)
	alertRepo := outrepo.NewAlertRepository(pool)

	// 2) Usecases
	devUC := usecase.NewDevicesUsecase(devRepo, planRepo, alertRepo)

	// 3) Handlers
	devH := handler.NewDevicesHandler(devUC)

	// 4) Router gốc (đã gắn Recovery, RequestID, Logger, CORS, Prometheus, healthz/readiness, /metrics)
	r := router.New(pool, baseLogger, router.Options{
		AppEnv:      cfg.AppEnv,
		AllowOrigin: cfg.AllowOrigin,
	})

	// 5) Mount modules vào /api
	api := r.Group("/api")
	router.MountDevices(api, devH)

	return r
}

// RunHTTP: chạy server với graceful shutdown
func RunHTTP(r *gin.Engine, port string) error {
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Start server
	errCh := make(chan error, 1)
	go func() {
		log.Printf("HTTP listening on :%s", port)
		errCh <- srv.ListenAndServe()
	}()

	// Wait signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		log.Printf("shutdown signal: %s", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
