-- Dev-only reset + seed data for local Postgres.

TRUNCATE TABLE messages, channels, users RESTART IDENTITY CASCADE;

INSERT INTO users (id, username, display_name)
VALUES
    (1, 'daniel', 'Daniel');

INSERT INTO channels (id, name)
VALUES
    (1, 'general');

INSERT INTO messages (author_id, channel_id, body)
VALUES
    (1, 1, 'Welcome to PulseLounge'),
    (1, 1, 'Seed data helps unblock UI and API testing'),
    (1, 1, 'Run make seed-reset-dev to clear + reseed');
