## 1. Định nghĩa struct chứa Handler
``` go
// Định nghĩa struct để "đóng gói" các handler
type RouterDeps struct {
    Devices *DevicesHandler    // Field chứa pointer đến DevicesHandler
    Users   *UsersHandler      // Field chứa pointer đến UsersHandler
}
```
## 2. Tạo instance của Handler
``` go
// Tạo các handler instance
devicesHandler := &DevicesHandler{}  // Tạo DevicesHandler
usersHandler := &UsersHandler{}      // Tạo UsersHandler

// Đóng gói vào struct
deps := RouterDeps{
    Devices: devicesHandler,    // Gán devicesHandler vào field Devices
    Users:   usersHandler,      // Gán usersHandler vào field Users
}
```
## 3. Function Router nhận struct
``` go
func NewRouter(d RouterDeps) *gin.Engine {
    //         ↑ ↑
    //         │ └── Kiểu RouterDeps (struct đã định nghĩa)
    //         └──── Tên parameter (tự đặt)
    
    router := gin.Default()
    
    // MÓKKẾT: d.Devices chính là devicesHandler đã truyền vào
    router.GET("/devices", d.Devices.GetDevices)
    //                     ↑        ↑
    //                     │        └── Method của DevicesHandler
    //                     └─────────── Field Devices từ struct RouterDeps
    
    return router
}
```
## 4. Sơ đồ móc nối:
``` 
devicesHandler (instance)
       │
       │ gán vào
       ▼
RouterDeps{
    Devices: devicesHandler  ←── Field "Devices"
}
       │
       │ truyền vào function
       ▼
func NewRouter(d RouterDeps)
       │
       │ sử dụng
       ▼
d.Devices.GetDevices
│    │         │
│    │         └── Method của DevicesHandler
│    └─────────── Field Devices của struct
└──────────────── Parameter d (kiểu RouterDeps)
```
## 5. Ví dụ đầy đủ:
``` go
// 1. Định nghĩa Handler
type DevicesHandler struct {
    // các field khác...
}

func (h *DevicesHandler) GetDevices(c *gin.Context) {
    // logic xử lý API
}

// 2. Định nghĩa struct đóng gói
type RouterDeps struct {
    Devices *DevicesHandler
}

// 3. Function Router
func NewRouter(d RouterDeps) *gin.Engine {
    router := gin.Default()
    
    // Móc nối: d.Devices trỏ đến instance DevicesHandler
    router.GET("/devices", d.Devices.GetDevices)
    
    return router
}

// 4. Sử dụng
func main() {
    // Tạo handler instance
    devicesH := &DevicesHandler{}
    
    // Đóng gói
    deps := RouterDeps{Devices: devicesH}
    
    // Truyền vào router
    r := NewRouter(deps)
    
    r.Run(":8080")
}
```
## 6. Chuỗi móc nối:
``` 
Instance → Field → Parameter → Usage
    ↓        ↓         ↓        ↓
devicesH → Devices → d → d.Devices.GetDevices
```
