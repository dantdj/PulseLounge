package channels

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

func (r PostgresRepository) List(ctx context.Context) (_ []Channel, err error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, created_at FROM channels ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	channels := make([]Channel, 0)
	for rows.Next() {
		var c Channel
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}

func (r PostgresRepository) Create(ctx context.Context, name string) (Channel, error) {
	var c Channel
	err := r.db.QueryRowContext(ctx, `
INSERT INTO channels (name)
VALUES ($1)
RETURNING id, name, created_at`, name).Scan(&c.ID, &c.Name, &c.CreatedAt)
	if err != nil {
		return Channel{}, err
	}
	return c, nil
}

func (r PostgresRepository) Delete(ctx context.Context, id int64) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, "DELETE FROM messages WHERE channel_id = $1", id); err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM channels WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrChannelNotFound
	}

	return tx.Commit()
}
