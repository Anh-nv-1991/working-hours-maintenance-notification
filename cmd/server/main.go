package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	mw "wh-ma/internal/adapter/inbound/http/middleware"
	"wh-ma/internal/bootstrap"
)

func main() {
	// 1) Context + signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2) Load config (đã có trong bootstrap)
	cfg := bootstrap.LoadConfig()

	// 3) Init OpenTelemetry (CẤY Ở ĐÂY)
	// Hàm này nằm ở internal/bootstrap/otel.go (anh thêm theo mẫu trước đó)
	tr := bootstrap.InitTracing(ctx, "wh-ma-api")
	defer func() { _ = tr.Shutdown(ctx) }()

	// 4) DB pool (nếu đã gắn tracer trong NewPGXPool thì mọi query sẽ có span)
	pool, err := bootstrap.NewPGXPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// 5) Router
	r := bootstrap.BuildRouter(cfg, pool, nil)

	// 6) OTel HTTP middleware của anh (đảm bảo có span cho mọi route)
	r.Use(mw.OTelMiddleware("wh-ma-api"))

	// 7) Run HTTP (graceful)
	if err := bootstrap.RunHTTP(r, cfg.Port); err != nil {
		log.Fatalf("http: %v", err)
	}
}
