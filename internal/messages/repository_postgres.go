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

func (r PostgresRepository) List(ctx context.Context) (_ []Message, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, body, created_at FROM messages ORDER BY id")
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
		if err := rows.Scan(&m.ID, &m.Body, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r PostgresRepository) Create(ctx context.Context, body string) (Message, error) {
	var m Message
	err := r.db.QueryRowContext(ctx, "INSERT INTO messages (body) VALUES ($1) RETURNING id, body, created_at", body).Scan(&m.ID, &m.Body, &m.CreatedAt)
	if err != nil {
		return Message{}, err
	}
	return m, nil
}

func (r PostgresRepository) Edit(ctx context.Context, id int, body string) error {
	row := r.db.QueryRowContext(ctx, "UPDATE messages SET body = $1 WHERE id = $2", body, id)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
