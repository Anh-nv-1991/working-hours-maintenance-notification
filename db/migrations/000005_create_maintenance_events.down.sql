-- 5_down
DROP INDEX IF EXISTS idx_maint_device_at;
ALTER TABLE maintenance_events DROP CONSTRAINT IF EXISTS fk_maint_device;
DROP TABLE IF EXISTS maintenance_events;
