package main

import (
	"context"
	"fmt"
	db "maint/internal/db/sqlc"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	// ưu tiên DATABASE_URL; nếu thiếu thì ráp từ DB_* trong .env

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	q := db.New(pool)

	dev, err := q.CreateDevice(ctx, "pc-2000")
	if err != nil {
		panic(err)
	}
	fmt.Println("inserted:", dev.ID, dev.Name)

	got, err := q.GetDevice(ctx, dev.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("fetched:", got.ID, got.Name)
}
