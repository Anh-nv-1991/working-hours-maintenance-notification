 ## 1. **Định nghĩa struct (khuôn mẫu)**

```go
type UserHandler struct {
    db     *sql.DB      
    logger *log.Logger  
}
```

**Giải thích:**

- `type UserHandler struct` = Định nghĩa kiểu dữ liệu mới tên "UserHandler" dạng struct
- `db *sql.DB` = Field tên "db", kiểu pointer đến sql.DB
- `logger *log.Logger` = Field tên "logger", kiểu pointer đến log.Logger
- `*` = Con trỏ (pointer), trỏ đến vùng nhớ chứa object

## 2. **Tạo instance (đối tượng thực)**

```go
userHandler := &UserHandler{
    db:     database,
    logger: logger,
}
```

**Giải thích từng phần:**

- `userHandler` = Tên biến
- `:=` = Khai báo và gán giá trị (short declaration)
- `&UserHandler` = Tạo instance mới của UserHandler và trả về pointer
- `{...}` = Khởi tạo giá trị cho các field
- `db: database` = Gán giá trị biến cho field `db` `database`
- `logger: logger` = Gán giá trị biến `logger` cho field `logger`

## 3. **Method của struct**

```go
func (h *UserHandler) GetUsers(c *gin.Context) {
    users := h.db.Query("SELECT * FROM users")
    h.logger.Info("Getting users")
}
```

**Giải thích:**

- `func` = Khai báo function

- `(h *UserHandler)` = **Receiver** - function này thuộc về struct UserHandler
  
  - `h` = Tên receiver (có thể đặt tên gì cũng được)
  - `*UserHandler` = Kiểu receiver (pointer đến UserHandler)

- `GetUsers` = Tên method

- `(c *gin.Context)` = Parameter của method

- `h.db` = Truy cập field `db` của instance thông qua receiver `h`

- `h.logger` = Truy cập field `logger` của instance thông qua receiver `h`

## 4. **Sử dụng instance trong router**

```go
router.GET("/users", userHandler.GetUsers)
```

**Giải thích:**

- `router.GET` = Gọi method GET của object router
- `"/users"` = URL path (string)
- `userHandler.GetUsers` = Truy cập method GetUsers của instance userHandler
- **KHÔNG có ()** = Chỉ truyền reference của method, không gọi method

## 5. **Sơ đồ cú pháp:**

```
                    Receiver
                       ↓
func (h *UserHandler) GetUsers(c *gin.Context) {
 ↑    ↑      ↑         ↑           ↑
 │    │      │         │           │
 │    │      │         │           └── Parameter type
 │    │      │         └────────────── Method name  
 │    │      └──────────────────────── Receiver type
 │    └─────────────────────────────── Receiver name
 └────────────────────────────────────── Function keyword
```

## 6. **Chi tiết về Receiver:**

```go
// Cú pháp receiver
func (receiverName *StructType) MethodName() {
//    ↑            ↑              ↑
//    │            │              └── Tên method
//    │            └─────────────── Kiểu struct  
//    └──────────────────────────── Tên receiver (tự đặt)
}

// Ví dụ cụ thể
func (h *UserHandler) GetUsers() {
//    ↑  ↑             ↑
//    │  │             └── Method name
//    │  └─────────────── Struct type
//    └────────────────── Receiver name (có thể đặt u, handler, gì cũng được)
}
```

## 7. **Ví dụ step-by-step:**

```go
// Bước 1: Tạo struct
type Calculator struct {
    result int
}

// Bước 2: Tạo instance
calc := &Calculator{result: 0}
//  ↑         ↑          ↑
//  │         │          └── Khởi tạo field
//  │         └─────────── Tạo instance với &
//  └───────────────────── Tên biến

// Bước 3: Định nghĩa method với receiver
func (c *Calculator) Add(num int) {
//    ↑  ↑            ↑      ↑
//    │  │            │      └── Parameter
//    │  │            └─────── Method name
//    │  └──────────────────── Receiver type
//    └─────────────────────── Receiver name
    c.result += num  // Dùng receiver để truy cập field
}

// Bước 4: Gọi method
calc.Add(5)  // Gọi method Add của instance calc
//  ↑   ↑  ↑
//  │   │  └── Argument
//  │   └───── Method name
//  └────────── Instance name
```
