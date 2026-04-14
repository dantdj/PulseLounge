-- Dev-only seed data for local Postgres
-- Safe to run multiple times.

INSERT INTO users (id, username, display_name, created_at)
VALUES
    (1, 'daniel', 'Daniel', NOW() - INTERVAL '4 hours')
ON CONFLICT (id) DO NOTHING;

INSERT INTO channels (id, name, created_at)
VALUES
    (1, 'general', NOW() - INTERVAL '4 hours')
ON CONFLICT (id) DO NOTHING;

INSERT INTO messages (id, author_id, channel_id, body, created_at)
VALUES
    (1, 1, 1, 'Welcome to PulseLounge', NOW() - INTERVAL '3 hours'),
    (2, 1, 1, 'Seed data helps unblock UI and API testing', NOW() - INTERVAL '2 hours'),
    (3, 1, 1, 'Run make seed-reset-dev to clear + reseed', NOW() - INTERVAL '1 hour')
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM users), 1), true);
SELECT setval('channels_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM channels), 1), true);

-- Keep sequence in sync when explicit IDs are used.
SELECT setval('messages_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM messages), 1), true);
