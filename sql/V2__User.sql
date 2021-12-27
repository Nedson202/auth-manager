BEGIN;

CREATE SCHEMA identity;

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE identity.user (
    id text PRIMARY KEY,
    username text UNIQUE,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TRIGGER user_update_trigger
AFTER UPDATE ON identity.user
FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_timestamp();

COMMIT;
