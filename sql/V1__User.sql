BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY,
    username text UNIQUE,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    role text NOT NULL DEFAULT 'User',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

COMMIT;
