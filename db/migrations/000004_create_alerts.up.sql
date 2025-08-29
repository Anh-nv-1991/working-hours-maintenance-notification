-- 4_up: create alerts
CREATE TABLE IF NOT EXISTS alerts (
  id          BIGSERIAL PRIMARY KEY,
  device_id   BIGINT NOT NULL,
  type        TEXT NOT NULL,
  message     TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  resolved    BOOLEAN NOT NULL DEFAULT FALSE,
  resolved_at TIMESTAMPTZ,
  resolved_by TEXT
);

ALTER TABLE alerts
  ADD CONSTRAINT fk_alerts_device
  FOREIGN KEY (device_id) REFERENCES devices(id)
  ON UPDATE CASCADE ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_alerts_device_created ON alerts(device_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_alerts_resolved ON alerts(resolved);
