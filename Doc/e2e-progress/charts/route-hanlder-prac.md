Bài tập 1 — “Hello route” (hiểu route đơn giản và handler)
- Mục tiêu: tạo server có 1 route GET /health hoặc /hello, tách handler ra file riêng.
- Cấu trúc đề xuất:
    - cmd/server/main.go
    - internal/handlers/health.go

- Nội dung mẫu:
``` go
package main

import (
	"fmt"
	"net/http"
)

// internal/handlers/health.go (nên đặt vào package handlers trong thực tế, 
// đây là ví dụ đơn giản gom vào cùng package main cho dễ chạy)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler) // route -> handler

	http.ListenAndServe(":8080", mux)
}
```
- Cách chạy: go run cmd/server/main.go
- Thử bằng curl:
    - curl -i [http://localhost:8080/health](http://localhost:8080/health)

- Kỳ vọng: HTTP 200 và body "OK"

Bài tập 2 — Route có path param và REST-ish (GET /users/{id}, POST /users)
- Mục tiêu: hiểu cách tách handler, xử lý path param, đọc JSON body cho POST.
- Cấu trúc đề xuất:
    - cmd/server/main.go
    - internal/handlers/users.go
    - internal/models/user.go (nếu muốn mở rộng)

- Nội dung mẫu (dùng standard net/http; lấy id từ URL bằng TrimPrefix):
``` go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Simulate in-memory store
var users = map[string]User{}

// Handler: POST /users  (create)
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if u.ID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	users[u.ID] = u
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// Handler: GET /users/{id}
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	prefix := "/users/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, prefix)
	if id == "" {
		http.NotFound(w, r)
		return
	}
	u, ok := users[id]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(u)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", createUserHandler) // POST /users
	mux.HandleFunc("/users/", getUserHandler)   // GET /users/{id}

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", mux)
}
```
- Thử:
    - Tạo user: curl -X POST -d '{"id":"u1","name":"An"}' -H "Content-Type: application/json" [http://localhost:8080/users](http://localhost:8080/users)
    - Lấy user: curl [http://localhost:8080/users/u1](http://localhost:8080/users/u1)

- Kỳ vọng: POST trả về 201 với JSON, GET trả về JSON user

Gợi ý mở rộng: tách handlers ra package internal/handlers, models vào internal/models; main chỉ khởi tạo router và register handlers.
Bài tập 3 — Tách router, handlers và middleware (cấu trúc rõ ràng hơn)
- Mục tiêu: nắm cách tổ chức code theo lớp: cmd (entry), internal/router (đăng ký routes), internal/handlers (logic), optional internal/middleware.
- Cấu trúc đề xuất:
    - cmd/server/main.go
    - internal/router/router.go
    - internal/handlers/handlers.go
    - internal/middleware/logging.go

- Ví dụ ngắn minh hoạ cách tách (chỉ ý tưởng — bạn có thể triển khai đầy đủ):
``` go
// cmd/server/main.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := NewRouter() // từ internal/router
	fmt.Println("Listening :8080")
	http.ListenAndServe(":8080", mux)
}
```

``` go
// internal/router/router.go
package main

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// đăng ký handler (từ package handlers)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/users", createUserHandler)
	mux.HandleFunc("/users/", getUserHandler)
	// nếu dùng middleware, bọc handler ở đây
	return mux
}
```

``` go
// internal/middleware/logging.go
package main

import (
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
```
