## Cấu trúc thư mục:

```
user-management/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   │   └── user.go
│   ├── usecase/
│   │   └── user_usecase.go
│   ├── adapter/
│   │   ├── inbound/
│   │   │   ├── http/
│   │   │   │   ├── handler.go
│   │   │   │   └── request.go
│   │   │   └── port/
│   │   │       └── user_service.go
│   │   └── outbound/
│   │       ├── repository/
│   │       │   └── user_repository.go
│   │       └── port/
│   │           └── user_repository.go
└── go.mod
mkdir user-management
cd user-management

mkdir cmd
mkdir internal

mkdir internal\domain
mkdir internal\usecase

mkdir internal\adapter
mkdir internal\adapter\inbound\http
mkdir internal\adapter\inbound\port
mkdir internal\adapter\outbound\repository
mkdir internal\adapter\outbound\port

# Tạo file rỗng
ni go.mod -ItemType File
ni cmd\main.go -ItemType File
ni internal\domain\user.go -ItemType File
ni internal\usecase\user_usecase.go -ItemType File
ni internal\adapter\inbound\http\handler.go -ItemType File
ni internal\adapter\inbound\http\request.go -ItemType File
ni internal\adapter\inbound\port\user_service.go -ItemType File
ni internal\adapter\outbound\repository\user_repository.go -ItemType File
ni internal\adapter\outbound\port\user_repository.go -ItemType File

```

## 1. Domain Layer

**internal/domain/user.go**

```go
package domain

import "time"

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

func NewUser(name string) *User {
    return &User{
        Name:      name,
        CreatedAt: time.Now(),
    }
}
```

## 2. Outbound Port (Interface)

**internal/adapter/outbound/port/user_repository.go**

```go
package port

import "user-management/internal/domain"

type UserRepository interface {
    Save(user *domain.User) error
    FindAll() ([]*domain.User, error)
    FindByID(id int) (*domain.User, error)
}
```

## 3. Outbound Adapter (Repository Implementation)

**internal/adapter/outbound/repository/user_repository.go**

```go
package repository

import (
    "fmt"
    "sync"
    "user-management/internal/adapter/outbound/port"
    "user-management/internal/domain"
)

type InMemoryUserRepository struct {
    users  map[int]*domain.User
    nextID int
    mutex  sync.RWMutex
}

func NewInMemoryUserRepository() port.UserRepository {
    return &InMemoryUserRepository{
        users:  make(map[int]*domain.User),
        nextID: 1,
    }
}

func (r *InMemoryUserRepository) Save(user *domain.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    user.ID = r.nextID
    r.users[r.nextID] = user
    r.nextID++

    return nil
}

func (r *InMemoryUserRepository) FindAll() ([]*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    users := make([]*domain.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }

    return users, nil
}

func (r *InMemoryUserRepository) FindByID(id int) (*domain.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    user, exists := r.users[id]
    if !exists {
        return nil, fmt.Errorf("user with id %d not found", id)
    }

    return user, nil
}
```

## 4. Inbound Port (Service Interface)

**internal/adapter/inbound/port/user_service.go**

```go
package port

import "user-management/internal/domain"

type UserService interface {
    CreateUser(name string) (*domain.User, error)
    GetAllUsers() ([]*domain.User, error)
    GetUserByID(id int) (*domain.User, error)
}
```

## 5. Use Case Layer

**internal/usecase/user_usecase.go**

```go
package usecase

import (
    "user-management/internal/adapter/inbound/port"
    outboundPort "user-management/internal/adapter/outbound/port"
    "user-management/internal/domain"
)

type UserUseCase struct {
    userRepo outboundPort.UserRepository
}

func NewUserUseCase(userRepo outboundPort.UserRepository) port.UserService {
    return &UserUseCase{
        userRepo: userRepo,
    }
}

func (uc *UserUseCase) CreateUser(name string) (*domain.User, error) {
    if name == "" {
        return nil, fmt.Errorf("name cannot be empty")
    }

    user := domain.NewUser(name)

    err := uc.userRepo.Save(user)
    if err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }

    return user, nil
}

func (uc *UserUseCase) GetAllUsers() ([]*domain.User, error) {
    users, err := uc.userRepo.FindAll()
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }

    return users, nil
}

func (uc *UserUseCase) GetUserByID(id int) (*domain.User, error) {
    user, err := uc.userRepo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return user, nil
}
```

## 6. Request Models

**internal/adapter/inbound/http/request.go**

```go
package http

type CreateUserRequest struct {
    Name string `json:"name" binding:"required"`
}

type GetUserRequest struct {
    ID int `json:"id" uri:"id" binding:"required"`
}
```

## 7. HTTP Handler

**internal/adapter/inbound/http/handler.go**

```go
package http

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "user-management/internal/adapter/inbound/port"
)

type UserHandler struct {
    userService port.UserService
}

func NewUserHandler(userService port.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.userService.CreateUser(req.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    user, err := h.userService.GetUserByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
    userGroup := r.Group("/api/users")
    {
        userGroup.POST("", h.CreateUser)
        userGroup.GET("", h.GetAllUsers)
        userGroup.GET("/:id", h.GetUserByID)
    }
}
```

## 8. Main Application

**cmd/main.go**

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "user-management/internal/adapter/inbound/http"
    "user-management/internal/adapter/outbound/repository"
    "user-management/internal/usecase"
)

func main() {
    // Khởi tạo repository (outbound adapter)
    userRepo := repository.NewInMemoryUserRepository()

    // Khởi tạo use case
    userService := usecase.NewUserUseCase(userRepo)

    // Khởi tạo handler (inbound adapter)
    userHandler := http.NewUserHandler(userService)

    // Khởi tạo Gin router
    r := gin.Default()

    // Đăng ký routes
    userHandler.RegisterRoutes(r)

    // Chạy server
    log.Println("Server starting on port 8080...")
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

## 9. Go Module

**go.mod**

```go
module user-management

go 1.21

require github.com/gin-gonic/gin v1.9.1
```

## Cách sử dụng:

1. **Tạo user mới:**
   
   ```bash
   curl -X POST http://localhost:8080/api/users \
   -H "Content-Type: application/json" \
   -d '{"name": "Nguyen Van A"}'
   ```
2. **Lấy tất cả users:**
   
   ```bash
   curl http://localhost:8080/api/users
   ```
3. **Lấy user theo ID:**
   
   ```bash
   curl http://localhost:8080/api/users/1
   ```
