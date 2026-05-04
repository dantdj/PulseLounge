import React from "react";
import { useMessages } from "../hooks/useMessages";
import type { Channel } from "../types";
import { MessageComposer } from "./MessageComposer";
import { MessageList } from "./MessageList";

type ChatViewProps = {
  selectedChannel: Channel | null;
  selectedChannelId: number | null;
};

export function ChatView({ selectedChannel, selectedChannelId }: ChatViewProps) {
  const [newMessageBody, setNewMessageBody] = React.useState("");
  const {
    messages,
    loading,
    error,
    submitError,
    isSubmitting,
    loadMessages,
    submitMessage,
    saveEditedMessage,
  } = useMessages(selectedChannelId);

  React.useEffect(() => {
    setNewMessageBody("");
  }, [selectedChannelId]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const submitted = await submitMessage(newMessageBody);
    if (submitted) {
      setNewMessageBody("");
    }
  };

  return (
    <section className="card">
      <div className="messages-header">
        <h2>{selectedChannel ? `#${selectedChannel.name}` : "Messages"}</h2>
      </div>

      <MessageList
        messages={messages}
        loading={loading}
        error={error}
        onEditMessage={saveEditedMessage}
      />

      <MessageComposer
        value={newMessageBody}
        isSubmitting={isSubmitting || selectedChannelId === null}
        submitError={submitError}
        onChange={setNewMessageBody}
        onRefresh={() => {
          void loadMessages();
        }}
        onSubmit={(event) => {
          void handleSubmit(event);
        }}
      />
    </section>
  );
}
