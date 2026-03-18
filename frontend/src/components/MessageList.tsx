import type { Message } from "../types";

type MessageListProps = {
  messages: Message[];
  loading: boolean;
  error: string | null;
};

export function MessageList({ messages, loading, error }: MessageListProps) {
  return (
    <div className="messages-scroll">
      {loading && <p>Loading messages...</p>}
      {error && <p className="error">Failed to load messages: {error}</p>}

      {!loading && !error && (
        <>
          {messages.length === 0 ? (
            <p className="empty-state">No messages yet.</p>
          ) : (
            <ul className="message-list">
              {messages.map((message) => (
                <li key={message.id} className="message-item">
                  <p className="message-body">{message.body}</p>
                  <p className="message-meta">
                    #{message.id} · {new Date(message.created_at).toLocaleString()}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </>
      )}
    </div>
  );
}
