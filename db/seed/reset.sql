-- Dev-only reset + seed data for local Postgres.

TRUNCATE TABLE messages RESTART IDENTITY;

INSERT INTO messages (body)
VALUES
    ('Welcome to PulseLounge'),
    ('Seed data helps unblock UI and API testing'),
    ('Run make seed-reset-dev to clear + reseed');
