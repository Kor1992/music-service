CREATE TABLE music.generation_quene(
    id SERIAL PRIMARY KEY,
    track_id INTEGER NOT NULL REFERENCES music.tracks(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_generation_quene_status ON music.generation_quene(status);