-- +goose Up
ALTER TABLE messages
ADD COLUMN author_id BIGINT NOT NULL,
ADD COLUMN channel_id BIGINT NOT NULL,
ADD COLUMN edited_at TIMESTAMPTZ;

ALTER TABLE messages
ADD CONSTRAINT messages_author_id_fkey
FOREIGN KEY (author_id) REFERENCES users (id);

ALTER TABLE messages
ADD CONSTRAINT messages_channel_id_fkey
FOREIGN KEY (channel_id) REFERENCES channels (id);

CREATE INDEX messages_author_id_idx
ON messages (author_id);

CREATE INDEX messages_channel_id_created_at_idx
ON messages (channel_id, created_at);

-- +goose Down
DROP INDEX messages_author_id_idx;
DROP INDEX messages_channel_id_created_at_idx;

ALTER TABLE messages
DROP CONSTRAINT messages_author_id_fkey;

ALTER TABLE messages
DROP CONSTRAINT messages_channel_id_fkey;

ALTER TABLE messages
DROP COLUMN edited_at,
DROP COLUMN author_id,
DROP COLUMN channel_id;
