-- 3_up: create readings
CREATE TABLE IF NOT EXISTS readings (
  id          BIGSERIAL PRIMARY KEY,
  device_id   BIGINT NOT NULL,
  at          TIMESTAMPTZ NOT NULL,
  hours_delta INTEGER  NOT NULL CHECK (hours_delta >= 0),
  location    TEXT,
  operator_id TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- FK tới devices
ALTER TABLE readings
  ADD CONSTRAINT fk_readings_device
  FOREIGN KEY (device_id) REFERENCES devices(id)
  ON UPDATE CASCADE ON DELETE CASCADE;

-- Chỉ mục hay dùng
CREATE INDEX IF NOT EXISTS idx_readings_device_at ON readings(device_id, at DESC);
