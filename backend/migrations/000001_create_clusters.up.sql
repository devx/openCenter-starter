CREATE TABLE IF NOT EXISTS clusters (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS clusters_status_idx ON clusters (status);
CREATE INDEX IF NOT EXISTS clusters_name_idx ON clusters (name);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS clusters_set_updated_at ON clusters;
CREATE TRIGGER clusters_set_updated_at
BEFORE UPDATE ON clusters
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
