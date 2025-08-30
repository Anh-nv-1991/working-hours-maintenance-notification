-- 4_down
DROP INDEX IF EXISTS idx_alerts_resolved;
DROP INDEX IF EXISTS idx_alerts_device_created;
ALTER TABLE alerts DROP CONSTRAINT IF EXISTS fk_alerts_device;
DROP TABLE IF EXISTS alerts;
