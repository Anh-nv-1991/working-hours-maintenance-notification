CREATE TABLE devices (
  id BIGSERIAL PRIMARY KEY,
  serial_number TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,

  -- profile
  model TEXT,
  manufacturer TEXT,
  year_of_manufacture INT,
  commission_date DATE,

  -- state
  total_working_hour INT DEFAULT 0,
  after_overhaul_working_hour INT DEFAULT 0,
  last_service_at TIMESTAMPTZ,
  location TEXT,
  avg_daily_hours FLOAT8,
  expected_next_maint TIMESTAMPTZ,

  -- status
  status TEXT NOT NULL DEFAULT 'active',

  -- audit
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  created_by TEXT,
  updated_by TEXT,
  deleted_by TEXT
);
