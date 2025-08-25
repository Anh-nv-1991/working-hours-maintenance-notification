// cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"

	"maint/internal/bootstrap"

	// compile-time sanity checks (giúp lệch chữ ký là fail build)
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

	// 5) Devices
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
			// Xử lý duplicate theo unique index (Postgres 23505)
			if strings.Contains(err.Error(), "SQLSTATE 23505") {
				c.JSON(http.StatusConflict, gin.H{"error": "device name already exists"})
				return
			}
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

	// 6) Readings
	// POST /readings { "device_id": 1, "value": 87.5, "at": "2025-08-25T10:00:00Z" }
	r.POST("/readings", func(c *gin.Context) {
		var req struct {
			DeviceID int64      `json:"device_id" binding:"required"`
			Value    float64    `json:"value" binding:"required"`
			At       *time.Time `json:"at,omitempty"` // tùy chọn
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var at interface{}
		if req.At != nil {
			at = *req.At
		} else {
			at = nil // để NOW()
		}
		reading, err := deps.Readings.AddReading(c, db.AddReadingParams{
			DeviceID: req.DeviceID,
			Value:    req.Value,
			Column3:  at,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, reading)
	})

	// GET /readings/last/:device_id
	r.GET("/readings/last/:device_id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
			return
		}
		reading, err := deps.Readings.LastReading(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, reading)
	})

	// 7) Plans
	// POST /plans { "device_id":1, "threshold_min":10, "threshold_max":100, "created_at": "..."? }
	r.POST("/plans", func(c *gin.Context) {
		var req struct {
			DeviceID     int64      `json:"device_id" binding:"required"`
			ThresholdMin float64    `json:"threshold_min" binding:"required"`
			ThresholdMax float64    `json:"threshold_max" binding:"required"`
			CreatedAt    *time.Time `json:"created_at,omitempty"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var created interface{}
		if req.CreatedAt != nil {
			created = *req.CreatedAt
		} else {
			created = nil // NOW()
		}
		plan, err := deps.Plans.UpsertPlan(c, db.UpsertPlanParams{
			DeviceID:     req.DeviceID,
			ThresholdMin: req.ThresholdMin,
			ThresholdMax: req.ThresholdMax,
			Column4:      created,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, plan)
	})

	// GET /plans/:device_id
	r.GET("/plans/:device_id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
			return
		}
		plan, err := deps.Plans.GetPlanForDevice(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, plan)
	})

	// 8) Alerts
	// POST /alerts/compute/:device_id  -> check breach & tạo alert nếu LOW/HIGH
	r.POST("/alerts/compute/:device_id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("device_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
			return
		}
		res, err := deps.Alerts.CheckThresholdBreach(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if res.Status == "OK" {
			c.JSON(http.StatusOK, gin.H{"status": "OK"})
			return
		}
		alert, err := deps.Alerts.CreateAlert(c, db.CreateAlertParams{
			DeviceID:  id,
			ReadingID: pgtype.Int8{Int64: res.ReadingID, Valid: true},
			Level:     res.Status,
			Message:   pgtype.Text{String: fmt.Sprintf("Value %.2f breached", res.ReadingValue), Valid: true},
			Column5:   nil, // NOW()
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, alert)
	})

	// POST /alerts/:id/service -> mark serviced
	r.POST("/alerts/:id/service", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		alert, err := deps.Alerts.ResolveAlert(c, db.ResolveAlertParams{ID: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, alert)
	})

	// 9) Start server + graceful shutdown
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
