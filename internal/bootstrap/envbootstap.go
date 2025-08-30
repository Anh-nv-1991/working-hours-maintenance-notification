package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

// ưu tiên configs/.env; nếu không có thì thử .env ở root.
// Overload để giá trị trong file .env ghi đè môi trường local khi dev.
// (Trong production, bạn nên để biến env từ hệ thống/compose và bỏ Overload)
func LoadEnvFirst() {
	// nếu đã có biến APP_ENV từ môi trường, vẫn cho phép .env ghi đè khi dev
	if _, err := os.Stat("configs/.env"); err == nil {
		_ = godotenv.Overload("configs/.env")
		return
	}
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Overload(".env")
	}
}
