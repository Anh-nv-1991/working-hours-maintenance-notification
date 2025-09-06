-- 3_down
DROP INDEX IF EXISTS idx_readings_device_at;
ALTER TABLE readings DROP CONSTRAINT IF EXISTS fk_readings_device;
DROP TABLE IF EXISTS readings;
