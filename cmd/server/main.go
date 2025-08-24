// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"maint/internal/bootstrap"

	// compile-time sanity checks (tùy chọn, an toàn)
	db "maint/internal/db/sqlc"
	"maint/internal/ports"
)

/*** Build sẽ FAIL nếu lệch chữ ký giữa sqlc và ports ***/
var (
	_ ports.DeviceRepo  = (*db.Queries)(nil)
	_ ports.ReadingRepo = (*db.Queries)(nil)
	_ ports.PlanRepo    = (*db.Queries)(nil)
	_ ports.AlertRepo   = (*db.Queries)(nil)
)

func main() {
	// 1) Env
	_ = godotenv.Load("configs/.env") // không lỗi nếu thiếu
	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2) Deps (DB pool + repos qua bootstrap)
	ctx := context.Background()
	deps, err := bootstrap.Init(ctx)
	if err != nil {
		log.Fatalf("bootstrap init: %v", err)
	}
	defer deps.Pool.Close()

	// 3) Router + middlewares
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	// 4) Health & Readiness
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true, "env": os.Getenv("APP_ENV")})
	})
	r.GET("/readiness", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := deps.Pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "db": "down", "err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "db": "up"})
	})

	// 5) Devices demo (khớp chữ ký sqlc hiện tại)
	// POST /devices { "name": "sensor-1" }
	r.POST("/devices", func(c *gin.Context) {
		var req struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		dev, err := deps.Devices.CreateDevice(c, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, dev)
	})

	// GET /devices/:id
	r.GET("/devices/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		dev, err := deps.Devices.GetDevice(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, dev)
	})

	// 6) Start server + graceful shutdown
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Printf("listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Wait for SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
}
