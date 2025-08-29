package bootstrap

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"wh-ma/internal/adapter/outbound/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// helper: parse int env
func getEnvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return def
}

// helper: parse duration env (ví dụ "300s", "1h")
func getEnvDuration(key string, def time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return def
}

// NewDB khởi tạo connection pool tới Postgres và return Queries (sqlc)
func NewDB() (*sqlc.Queries, *pgxpool.Pool) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("failed to parse DATABASE_URL: %v", err)
	}

	// đọc từ env
	cfg.MaxConns = int32(getEnvInt("DB_MAX_OPEN_CONNS", 10))
	cfg.MinConns = int32(getEnvInt("DB_MAX_IDLE_CONNS", 2))
	cfg.MaxConnLifetime = getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour)
	cfg.MaxConnIdleTime = getEnvDuration("DB_CONN_MAX_IDLE", 30*time.Minute)
	cfg.HealthCheckPeriod = getEnvDuration("DB_HEALTHCHECK_PERIOD", 30*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create pgx pool: %v", err)
	}

	// kiểm tra kết nối
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	log.Printf("✅ Connected to Postgres (MaxConns=%d, MinConns=%d)", cfg.MaxConns, cfg.MinConns)

	queries := sqlc.New(pool)
	return queries, pool
}
