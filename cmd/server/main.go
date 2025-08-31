package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wh-ma/internal/adapter/inbound/http/middleware"
	"wh-ma/internal/bootstrap"
)

func main() {
	cfg := bootstrap.LoadConfig()
	if cfg.DatabaseURL == "" {
		panic("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := bootstrap.NewPGXPool(ctx, cfg.DatabaseURL)
	if err != nil {
		panic("db connect failed: " + err.Error())
	}
	defer pool.Close()

	// 1) Tạo logger JSON theo LOG_LEVEL
	baseLogger := middleware.NewBaseLogger()

	// 2) Build router với middleware logger
	r := bootstrap.BuildRouter(cfg, pool, baseLogger)

	// 3) Chạy server trong goroutine
	errCh := make(chan error, 1)
	go func() { errCh <- bootstrap.RunHTTP(r, cfg.Port) }()

	// 4) Bắt tín hiệu dừng
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		baseLogger.Info("server.shutdown.signal", "sig", sig)
	case err := <-errCh:
		baseLogger.Error("server.http.error", "err", err)
	}

	// grace sleep
	time.Sleep(300 * time.Millisecond)
}
