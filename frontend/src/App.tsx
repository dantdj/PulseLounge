import React from "react";
import { MessageComposer } from "./components/MessageComposer";
import { MessageList } from "./components/MessageList";
import { useMessages } from "./hooks/useMessages";

export function App() {
  const [newMessageBody, setNewMessageBody] = React.useState("");
  const {
    messages,
    loading,
    error,
    submitError,
    isSubmitting,
    loadMessages,
    submitMessage,
  } = useMessages();

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const submitted = await submitMessage(newMessageBody);
    if (submitted) {
      setNewMessageBody("");
    }
  };

  return (
    <main className="app">
      <h1>PulseLounge</h1>

      <section className="card">
        <h2>Messages</h2>

        <MessageList messages={messages} loading={loading} error={error} />

        <MessageComposer
          value={newMessageBody}
          isSubmitting={isSubmitting}
          submitError={submitError}
          onChange={setNewMessageBody}
          onRefresh={loadMessages}
          onSubmit={handleSubmit}
        />
      </section>
    </main>
  );
}
