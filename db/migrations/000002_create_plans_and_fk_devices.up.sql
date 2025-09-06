-- 2_up: create plans + add FK to devices.plan_id (an toàn khi devices đã tồn tại)

CREATE TABLE IF NOT EXISTS plans (
  id               BIGSERIAL PRIMARY KEY,
  name             TEXT NOT NULL,
  interval_hours   INTEGER NOT NULL CHECK (interval_hours > 0),
  description      TEXT,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Thêm cột plan_id vào devices nếu chưa có
ALTER TABLE devices
  ADD COLUMN IF NOT EXISTS plan_id BIGINT NULL;

-- Thêm chỉ mục cho plan_id (nếu cần join/filter)
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_class c
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE c.relname = 'idx_devices_plan_id' AND n.nspname = 'public'
  ) THEN
    CREATE INDEX idx_devices_plan_id ON devices(plan_id);
  END IF;
END $$;

-- Thêm khóa ngoại nếu chưa có
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'fk_devices_plan_id'
  ) THEN
    ALTER TABLE devices
      ADD CONSTRAINT fk_devices_plan_id
      FOREIGN KEY (plan_id) REFERENCES plans(id)
      ON UPDATE CASCADE ON DELETE SET NULL;
  END IF;
END $$;
