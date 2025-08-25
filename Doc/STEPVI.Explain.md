1. Domain

Chứa các struct/entity: Device, Reading, Plan, Alert.

Đây là định nghĩa dữ liệu chuẩn mà toàn hệ thống dùng chung.

Không biết DB, không biết HTTP, chỉ là struct thuần Go.

Ví dụ:

type Reading struct {
    ID        int64
    DeviceID  int64
    Value     float64
    CreatedAt time.Time
}

2. Ports (interfaces)

Đặt trong internal/ports/.

Là cổng giao tiếp để use-case nói chuyện với “thế giới bên ngoài” (DB, clock, notifier).

Ví dụ:

type Repository interface {
    AddReading(ctx context.Context, deviceID int64, value float64, atUnix int64) (*domain.Reading, error)
}


👉 Use-case chỉ biết gọi Repository.AddReading() nhưng không biết DB Postgres hay Mongo.

3. Use-case

Là nghiệp vụ chính (AddReading, ComputeAlert, …).

Nhận input từ handler, dùng ports.Repository để đọc/ghi dữ liệu dạng domain.*.

Trả về domain.Reading, domain.Alert cho handler.

Ví dụ:

rd, err := uc.repo.AddReading(ctx, in.DeviceID, in.Value, uc.clock.NowUnix())
// rd là domain.Reading


👉 Use-case chỉ chơi với Domain + Ports.

4. Bootstrap

Chính là chỗ khâu nối dây.

Nó cầm lấy implementation thật (Repo dùng Postgres, Clock hệ thống, Notifier Noop/Slack) rồi gắn vào Use-case qua interface.

Ví dụ:

repo := repo.NewPGRepository(pool)   // implement ports.Repository
clk := clock.SystemClock{}           // implement ports.Clock
ntf := notifier.Noop{}               // implement ports.Notifier

addReadingUC := usecase.NewAddReadingUseCase(repo, clk, ntf)


👉 Bootstrap không xử lý dữ liệu, chỉ wiring.

5. main.go

Gọi bootstrap để build app.

Start server Gin và map endpoint.

Ví dụ:

readingsHandler := handlers.NewReadingsHandler(addReadingUC)
r.POST("/readings", readingsHandler.PostReading)

6. Quan hệ với API

Client → API call → Gin handler → Use-case → Ports → Repo → DB

Dữ liệu chạy qua các layer như sau:

(Client JSON)
   ↓
Gin Handler (parse JSON thành struct input)
   ↓
Use-case (xử lý logic, dùng Ports)
   ↓
Repo implement (query DB, map thành domain struct)
   ↓
Domain struct (Reading, Alert, …)
   ↓
Use-case trả ra
   ↓
Handler serialize thành JSON
   ↓
(API Response cho Client)

7. Tóm gọn quan hệ

Domain = model chuẩn.

Ports = interface mà Use-case cần.

Use-case = logic nghiệp vụ, chơi với Domain qua Ports.

Repo/Clock/Notifier = implement thật của Ports.

Bootstrap (main) = nối Use-case với Repo/Clock/Notifier thật.

Handler = cầu nối API ↔ Use-case.

👉 Hình dung:

Domain là “ngôn ngữ chung” (từ điển).

Ports là “ổ cắm điện” (interface).

Use-case là “thiết bị” (cần ổ cắm để chạy).

Bootstrap là “dây nối điện” (gắn ổ cắm với nguồn điện thật).

Main.go là “công tắc bật nguồn” (chạy hệ thống).