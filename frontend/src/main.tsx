import React from "react";
import ReactDOM from "react-dom/client";
import "./styles.css";

type Message = {
  id: number;
  body: string;
  created_at: string;
};

async function getErrorMessage(response: Response): Promise<string> {
  try {
    const data = (await response.json()) as { error?: string };
    if (data.error) {
      return data.error;
    }
  } catch {
    // Fall back to the status code when the response isn't JSON.
  }

  return `request failed with status ${response.status}`;
}

function App() {
  const [messages, setMessages] = React.useState<Message[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [newMessageBody, setNewMessageBody] = React.useState("");
  const [submitError, setSubmitError] = React.useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  const loadMessages = React.useCallback(() => {
    fetch("/api/messages")
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(await getErrorMessage(response));
        }
        return (await response.json()) as Message[];
      })
      .then((data) => {
        setMessages(data);
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : "unknown error";
        setError(message);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  React.useEffect(() => {
    loadMessages();
  }, [loadMessages]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const body = newMessageBody.trim();
    if (!body) {
      setSubmitError("Message cannot be empty.");
      return;
    }

    setSubmitError(null);
    setIsSubmitting(true);

    try {
      const response = await fetch("/api/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ body }),
      });
      if (!response.ok) {
        throw new Error(await getErrorMessage(response));
      }

      const created = (await response.json()) as Message;
      setMessages((current) => [...current, created]);
      setNewMessageBody("");
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "unknown error";
      setSubmitError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <main className="app">
      <h1>PulseLounge</h1>

      {loading && <p>Loading messages...</p>}
      {error && <p className="error">Failed to load messages: {error}</p>}

      {!loading && !error && (
        <section className="card">
          <h2>Messages</h2>

          <div className="messages-scroll">
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
          </div>

          <form className="message-form" onSubmit={handleSubmit}>
            <label htmlFor="message-body">New message</label>
            <textarea
              id="message-body"
              value={newMessageBody}
              onChange={(event) => setNewMessageBody(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter" && !event.shiftKey) {
                  event.preventDefault();
                  event.currentTarget.form?.requestSubmit();
                }
              }}
              rows={3}
              placeholder="Type a message..."
              disabled={isSubmitting}
            />
            <div className="message-form-actions">
              <button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Sending..." : "Send"}
              </button>
              <button type="button" onClick={loadMessages} disabled={isSubmitting}>
                Refresh
              </button>
            </div>
            {submitError && <p className="error">Failed to send message: {submitError}</p>}
          </form>
        </section>
      )}
    </main>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
