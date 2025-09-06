-- 2_down: gỡ FK + cột + bảng
ALTER TABLE devices DROP CONSTRAINT IF EXISTS fk_devices_plan_id;
DROP INDEX IF EXISTS idx_devices_plan_id;
ALTER TABLE devices DROP COLUMN IF EXISTS plan_id;

DROP TABLE IF EXISTS plans;
