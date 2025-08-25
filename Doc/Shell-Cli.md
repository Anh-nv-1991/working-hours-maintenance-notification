Server-start
go run .\cmd\server
Server-ping
// step 5 cli shell test
# Health & readiness
Invoke-RestMethod http://localhost:8080/healthz
Invoke-RestMethod http://localhost:8080/readiness

# Devices
$body = @{ name = "sensor-" + (Get-Random) } | ConvertTo-Json
Invoke-RestMethod -Method POST http://localhost:8080/devices -ContentType 'application/json' -Body $body
Invoke-RestMethod http://localhost:8080/devices/1

// limit to trigger alert
# Plans
$plan = @{ device_id=1; threshold_min=10; threshold_max=50 } | ConvertTo-Json
Invoke-RestMethod -Method POST http://localhost:8080/plans -ContentType 'application/json' -Body $plan
Invoke-RestMethod http://localhost:8080/plans/1

// read input data - working hours
# Readings
$reading = @{ device_id=1; value=60 } | ConvertTo-Json
Invoke-RestMethod -Method POST http://localhost:8080/readings -ContentType 'application/json' -Body $reading
Invoke-RestMethod http://localhost:8080/readings/last/1

// Show notification for users
# Alerts
Invoke-RestMethod -Method POST http://localhost:8080/alerts/compute/1
Invoke-RestMethod -Method POST http://localhost:8080/alerts/1/service

2) Attach vào Container (từ VS Code)

Cách nhanh (terminal):

# tên container xem bằng: docker ps
docker exec -it working_hours_maintenance_notification-api-1 sh
docker exec -it working_hours_maintenance_notification-db-1 psql -U postgres