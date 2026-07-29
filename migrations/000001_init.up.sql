CREATE SCHEMA IF NOT EXISTS music;

CREATE TABLE music.tracks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    prompt TEXT NOT NULL,
    audio_url TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    user_id INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_tracks_user_id ON music.tracks(user_id);
CREATE INDEX idx_tracks_status ON music.tracks(status);
