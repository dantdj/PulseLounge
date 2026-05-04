package messages

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PostgresRepository {
	return PostgresRepository{db: db}
}

func (r PostgresRepository) ListByChannel(ctx context.Context, channelID int64) (_ []Message, err error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, author_id, channel_id, body, created_at, edited_at
FROM messages
WHERE channel_id = $1
ORDER BY created_at, id`, channelID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	messages := make([]Message, 0)
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.AuthorID, &m.ChannelID, &m.Body, &m.CreatedAt, &m.EditedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r PostgresRepository) CreateInChannel(ctx context.Context, channelID int64, authorID int64, body string) (Message, error) {
	var m Message
	err := r.db.QueryRowContext(ctx, `
INSERT INTO messages (author_id, channel_id, body)
VALUES ($1, $2, $3)
RETURNING id, author_id, channel_id, body, created_at, edited_at`, authorID, channelID, body).
		Scan(&m.ID, &m.AuthorID, &m.ChannelID, &m.Body, &m.CreatedAt, &m.EditedAt)
	if err != nil {
		return Message{}, err
	}
	return m, nil
}

func (r PostgresRepository) Edit(ctx context.Context, id int64, body string) error {
	result, err := r.db.ExecContext(ctx, "UPDATE messages SET body = $1, edited_at = NOW() WHERE id = $2", body, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMessageNotFound
	}

	return nil
}
