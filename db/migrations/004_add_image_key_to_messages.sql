-- +goose Up
ALTER TABLE messages
ADD COLUMN image_key TEXT;

-- +goose Down
ALTER TABLE messages
DROP COLUMN image_key;