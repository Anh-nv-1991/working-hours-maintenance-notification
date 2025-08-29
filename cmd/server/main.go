// cmd/server/main.go
package main

import (
	"log"
	"os"
	healthrouter "wh-ma/internal/adapter/inbound/http/router"
	"wh-ma/internal/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	// Load env (đổi path nếu .env của ACE nằm trong configs/env/.env)
	_ = godotenv.Load("configs/.env")

	q, pool := bootstrap.NewDB() // q dùng cho các router khác (devices,…)
	_ = q
	defer pool.Close()

	r := healthrouter.New(pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
