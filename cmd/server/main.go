package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"maint/internal/bootstrap"
	"maint/internal/http/handlers"
)

func main() {
	// 1) Load env
	_ = godotenv.Load("configs/.env")
	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	port := getenv("PORT", "8080")

	// 2) Bootstrap (DB pool + sqlc repos)
	ctx := context.Background()
	deps, err := bootstrap.Init(ctx)
	if err != nil {
		log.Fatalf("bootstrap init: %v", err)
	}
	defer deps.Pool.Close()

	// 3) Init handlers (dùng repos từ deps)
	devicesH := handlers.NewDevicesHandler(deps.Devices)
	readingsH := handlers.NewReadingsHandler(deps.Readings)
	plansH := handlers.NewPlansHandler(deps.Plans)
	alertsH := handlers.NewAlertsHandler(deps.Alerts)

	// 4) Router & middlewares
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	// Health & readiness
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true, "env": os.Getenv("APP_ENV")})
	})
	r.GET("/readiness", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := deps.Pool.Ping(ctx); err != nil {
			c.JSON(503, gin.H{"ok": false, "db": "down", "err": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true, "db": "up"})
	})

	// 5) Routes (map vào handlers)
	// Devices
	r.POST("/devices", devicesH.PostDevice)
	r.GET("/devices/:id", devicesH.GetDevice)

	// Readings
	r.POST("/readings", readingsH.PostReading)
	r.GET("/readings/last/:device_id", readingsH.GetLastReading)

	// Plans
	r.POST("/plans", plansH.UpsertPlan)
	r.GET("/plans/:device_id", plansH.GetPlan)

	// Alerts
	r.POST("/alerts/compute/:device_id", alertsH.ComputeAlert)
	r.POST("/alerts/:id/service", alertsH.MarkServiced)

	// 6) Run server
	log.Printf("listening on :%s ...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
