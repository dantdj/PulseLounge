package messages

import "time"

type Message struct {
	ID        int64      `json:"id"`
	AuthorID  int64      `json:"author_id"`
	ChannelID int64      `json:"channel_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	EditedAt  *time.Time `json:"edited_at"`
}
