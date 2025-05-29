CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email TEXT NOT NULL UNIQUE,
                       name TEXT NOT NULL,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       modified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       last_login TIMESTAMPTZ,
                       is_enabled BOOLEAN NOT NULL DEFAULT TRUE
);

-- Automatically update `modified_at` on UPDATE
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.modified_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER set_modified_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
