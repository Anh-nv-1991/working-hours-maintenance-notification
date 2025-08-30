-- 5_up: create maintenance_events (lịch sử bảo dưỡng/sửa chữa)
CREATE TABLE IF NOT EXISTS maintenance_events (
  id           BIGSERIAL PRIMARY KEY,
  device_id    BIGINT NOT NULL,
  at           TIMESTAMPTZ NOT NULL,
  interval     INTEGER,            -- Khoảng bảo dưỡng áp dụng (giờ), có thể NULL nếu không theo chính sách
  notes        TEXT,
  performed_by TEXT,
  cost         NUMERIC(12,2) DEFAULT 0 CHECK (cost >= 0),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE maintenance_events
  ADD CONSTRAINT fk_maint_device
  FOREIGN KEY (device_id) REFERENCES devices(id)
  ON UPDATE CASCADE ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_maint_device_at ON maintenance_events(device_id, at DESC);
