CREATE TABLE music.users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    subscription TEXT NOT NULL DEFAULT 'trial',
    trial_ends_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '14 days',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);