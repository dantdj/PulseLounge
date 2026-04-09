import { EditableMessageItem } from "./EditableMessageItem";
import type { Message } from "../types";

type MessageListProps = {
  messages: Message[];
  loading: boolean;
  error: string | null;
  onEditMessage: (id: number, body: string) => Promise<string | null>;
};

export function MessageList({ messages, loading, error, onEditMessage }: MessageListProps) {
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
                <EditableMessageItem
                  key={message.id}
                  message={message}
                  onEditMessage={onEditMessage}
                />
              ))}
            </ul>
          )}
        </>
      )}
    </div>
  );
}
