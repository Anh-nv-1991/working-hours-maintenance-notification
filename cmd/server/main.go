package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wh-ma/internal/bootstrap"
)

func main() {
	cfg := bootstrap.LoadConfig()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := bootstrap.NewPGXPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	r := bootstrap.BuildRouter(cfg, pool)

	// Start server in goroutine for graceful shutdown
	errCh := make(chan error, 1)
	go func() { errCh <- bootstrap.RunHTTP(r, cfg.Port) }()

	// Wait for signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal: %v, shutting down...", sig)
	case err := <-errCh:
		log.Fatalf("http server error: %v", err)
	}

	// Grace period (if you implement http.Server with Shutdown, add here)
	time.Sleep(300 * time.Millisecond)
}
