# Kiểm tra cài đặt
docker --version
docker compose version

# Đăng chạy/ rebuild project
docker compose up -d --build

# Xem container đang chạy
docker ps

# Xem log theo service
docker compose logs -f api
docker compose logs -f db

# Dừng / xóa
docker compose stop
docker compose down        # xóa network + container
docker compose down -v     # xóa luôn volumes (cẩn thận mất DB)

# Rebuild nhanh khi đổi code
docker compose build api && docker compose up -d api

# Dọn dẹp rác
docker image prune -f
docker builder prune -f

##Lần đầu:

$env:DOCKER_BUILDKIT=1
docker compose up -d db
docker compose up --build -d api


##Các lần sau (đổi code):

docker compose up -d api


##Khi thay go.mod/go.sum hoặc Dockerfile:

docker compose build api
docker compose up -d api

//============Khởi động dự án trên local
🔹 1. Chạy host local (máy ACE trực tiếp, không Docker API)

👉 Dùng khi ACE đang dev code Go, muốn chạy nhanh.

Bước 1: Chạy Postgres bằng Docker (nếu chưa có DB local)
docker compose up db

Bước 2: Chạy API trực tiếp bằng Go CLI
go run ./cmd/server


📌 Trường hợp cần migrate DB trước:

migrate -path db/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up


🔹 2. Chạy toàn bộ trong Docker (API + DB)

👉 Dùng khi muốn test như production hoặc deploy.

Bước 1: Build lại image (nếu code thay đổi)
docker compose build

Bước 2: Chạy API + DB
docker compose up


📌 Thêm -d nếu muốn chạy background:

docker compose up -d