-- Dev-only seed data for local Postgres
-- Safe to run multiple times.

CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    body TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO messages (id, body, created_at)
VALUES
    (1, 'Welcome to PulseLounge', NOW() - INTERVAL '3 hours'),
    (2, 'Seed data helps unblock UI and API testing', NOW() - INTERVAL '2 hours'),
    (3, 'Run make seed-reset-dev to clear + reseed', NOW() - INTERVAL '1 hour')
ON CONFLICT (id) DO NOTHING;

-- Keep sequence in sync when explicit IDs are used.
SELECT setval('messages_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM messages), 1), true);
