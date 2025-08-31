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
