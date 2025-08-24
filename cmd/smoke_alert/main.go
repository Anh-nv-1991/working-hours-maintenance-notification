package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	db "maint/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := db.New(pool)

	// 1) Upsert plan (không có CreatedAt trong params)
	_, err = q.UpsertPlan(ctx, db.UpsertPlanParams{
		DeviceID:     1,
		ThresholdMin: 10,
		ThresholdMax: 90,
	})
	if err != nil {
		log.Fatal("upsert plan:", err)
	}

	// 2) Add reading (không có At trong params)
	rd, err := q.AddReading(ctx, db.AddReadingParams{
		DeviceID: 1,
		Value:    95,
	})
	if err != nil {
		log.Fatal("add reading:", err)
	}

	// 3) Check breach → tạo alert nếu cần
	chk, err := q.CheckThresholdBreach(ctx, 1)
	if err != nil {
		log.Fatal("check:", err)
	}

	if chk.Status != "OK" {
		al, err := q.CreateAlert(ctx, db.CreateAlertParams{
			DeviceID: 1,
			// ReadingID là pgtype.Int8 (nullable)
			ReadingID: pgtype.Int8{Int64: rd.ID, Valid: true},
			Level:     chk.Status,
			// Message là pgtype.Text (nullable)
			Message: pgtype.Text{
				String: fmt.Sprintf("Value %.2f out of range", chk.ReadingValue),
				Valid:  true,
			},
		})
		if err != nil {
			log.Fatal("create alert:", err)
		}
		fmt.Println("alert:", al.ID, al.Level)
	} else {
		fmt.Println("status:", chk.Status)
	}
}
