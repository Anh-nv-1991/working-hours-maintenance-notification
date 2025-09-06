# Domain Entities & Responsibilities

## Device (root aggregate)

- **Nhiệm vụ:** Đại diện cho một máy/thiết bị cụ thể ở mỏ.

- **Chứa:**
  
  - ID, SerialNumber, Name
  
  - Cấu hình tĩnh (**Profile**)
  
  - Trạng thái động (**State**)
  
  - Bộ đếm bảo dưỡng (**Counters**)
  
  - Status
  
  - Liên kết **PlanID**
  
  - Dấu vết thời gian/audit

- **Dùng khi:** CRUD thiết bị, cập nhật giờ vận hành, gán/bỏ kế hoạch bảo dưỡng (Plan), đổi trạng thái (active/maintenance/repair…).

---

## DeviceProfile (tĩnh)

- **Nhiệm vụ:** Hồ sơ “khai sinh” – model, hãng, năm sản xuất, ngày đưa vào khai thác.

- **Dùng khi:** hiển thị thông tin, phân tích vòng đời, đối chiếu khấu hao, lọc theo model/nhà sản xuất.

- **Ai cập nhật:** Thường chỉ khi khởi tạo thiết bị (ít thay đổi).

---

## OperationalState (động)

- **Nhiệm vụ:** Trạng thái vận hành hiện tại.

- **Thuộc tính:**
  
  - `Location`: vị trí hoạt động hiện tại
  
  - `TotalHours`: tổng giờ lũy kế (TWH)
  
  - `AfterOverhaul`: giờ sau đại tu (AOH) – tính chu kỳ bảo dưỡng từ mốc đại tu
  
  - `LastReadingAt`: mốc giờ vận hành ghi gần nhất
  
  - `ExpectedNextMaint`: dự đoán lần bảo dưỡng sắp tới (dựa trên Plan + tốc độ dùng)
  
  - `AvgDailyHours`: trung bình giờ/ngày (nền tảng dự báo)

- **Dùng khi:** cập nhật theo mỗi **Reading**, tính cảnh báo, dự báo lịch bảo dưỡng.

---

## MaintenancePolicy

- **Nhiệm vụ:** Quy tắc/định mức bảo dưỡng (ví dụ 250h thay dầu, 500h kiểm tra phanh…).

- **Dùng khi:** thiết kế policy chung theo loại máy; là đầu vào để tính **Counters** và cảnh báo.

---

## MaintenanceCounters & Counter (động)

- **Nhiệm vụ:** Theo dõi đã đạt mốc bảo dưỡng bao nhiêu lần với từng interval (key = IntervalHours).

- **Thuộc tính:**
  
  - `Counter.Count`: số lần đã làm
  
  - `Counter.LastAt`: mốc thời gian làm gần nhất
  
  - `Counter.Policy`: policy áp cho counter đó

- **Dùng khi:** tính đến hạn (due/overdue), tạo **Alert**, reset sau **MaintenanceEvent**.

---

## AuditMeta

- **Nhiệm vụ:** Lưu ai tạo/cập nhật/xóa – phục vụ traceability.

- **Dùng khi:** kiểm soát nội bộ, compliance, điều tra thay đổi.

---

## Reading (event giờ vận hành)

- **Nhiệm vụ:** Bản ghi “giờ vận hành tăng thêm” theo thời gian.

- **Thuộc tính:**
  
  - `HoursDelta`: số giờ tăng (đơn vị ca/ngày/tùy nhập liệu)
  
  - `At`, `Location`, `OperatorID`: bối cảnh giờ vận hành

- **Dùng khi:** tính lũy kế `TotalHours`, cập nhật `AvgDailyHours`, dự báo `ExpectedNextMaint`, kích hoạt kiểm tra cảnh báo.

---

## MaintenanceEvent (event bảo dưỡng/thực hiện)

- **Nhiệm vụ:** Ghi nhận đã làm bảo dưỡng tại mốc `Interval`.

- **Tác động:**
  
  - Cập nhật **Counters** (tăng Count, set LastAt, reset over-interval)
  
  - Ghi `Notes`, `PerformedBy`, `Cost` để truy vết và thống kê chi phí

- **Dùng khi:** sau khi xưởng/bộ phận bảo trì thực hiện xong.

---

## Alert (cảnh báo)

- **Nhiệm vụ:** Thông điệp hệ thống tạo ra để nhắc nhở/điều phối.

- **Thuộc tính:**
  
  - `Type`: `maintenance_due`, `over_usage`, `idle_too_long`…
  
  - `Resolved`: đánh dấu đã xử lý (ví dụ sau khi làm bảo dưỡng)

- **Dùng khi:** hiển thị dashboard/notification, SLA cảnh báo.
