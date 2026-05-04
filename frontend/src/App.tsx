import React from "react";
import { createChannel, deleteChannel, listChannels } from "./api/channels";
import { ChannelPicker } from "./components/ChannelPicker";
import { MessageComposer } from "./components/MessageComposer";
import { MessageList } from "./components/MessageList";
import { useMessages } from "./hooks/useMessages";
import type { Channel } from "./types";

function toErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "unknown error";
}

export function App() {
  const [newMessageBody, setNewMessageBody] = React.useState("");
  const [channels, setChannels] = React.useState<Channel[]>([]);
  const [selectedChannelId, setSelectedChannelId] = React.useState<number | null>(null);
  const [channelsLoading, setChannelsLoading] = React.useState(true);
  const [channelsError, setChannelsError] = React.useState<string | null>(null);
  const [isCreatingChannel, setIsCreatingChannel] = React.useState(false);
  const [newChannelName, setNewChannelName] = React.useState("");
  const [createChannelError, setCreateChannelError] = React.useState<string | null>(null);
  const [deleteChannelError, setDeleteChannelError] = React.useState<string | null>(null);
  const [channelCreateOpen, setChannelCreateOpen] = React.useState(false);
  const [deletingChannelId, setDeletingChannelId] = React.useState<number | null>(null);
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
    let cancelled = false;

    async function loadChannels() {
      setChannelsLoading(true);
      setChannelsError(null);

      try {
        const result = await listChannels();
        if (cancelled) {
          return;
        }

        setChannels(result);
        setSelectedChannelId((currentChannelId) => currentChannelId ?? result[0]?.id ?? null);
      } catch (loadChannelsError: unknown) {
        if (!cancelled) {
          setChannelsError(toErrorMessage(loadChannelsError));
        }
      } finally {
        if (!cancelled) {
          setChannelsLoading(false);
        }
      }
    }

    void loadChannels();

    return () => {
      cancelled = true;
    };
  }, []);

  const selectedChannel = channels.find((channel) => channel.id === selectedChannelId) ?? null;

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

  const handleDeleteChannel = async (channelId: number) => {
    setDeletingChannelId(channelId);
    setDeleteChannelError(null);

    try {
      await deleteChannel(channelId);
      setChannels((currentChannels) => {
        const nextChannels = currentChannels.filter((channel) => channel.id !== channelId);
        if (selectedChannelId === channelId) {
          setSelectedChannelId(nextChannels[0]?.id ?? null);
        }
        return nextChannels;
      });
    } catch (deleteChannelFailure: unknown) {
      setDeleteChannelError(toErrorMessage(deleteChannelFailure));
    } finally {
      setDeletingChannelId(null);
    }
  };

  const handleSubmitNewChannel = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const trimmedName = newChannelName.trim();
    if (!trimmedName) {
      setCreateChannelError("Channel name cannot be empty.");
      return;
    }

    setIsCreatingChannel(true);
    setCreateChannelError(null);

    try {
      const created = await createChannel(trimmedName);
      setChannels((currentChannels) => [...currentChannels, created]);
      setSelectedChannelId(created.id);
      setNewChannelName("");
      setChannelCreateOpen(false);
    } catch (createChannelFailure: unknown) {
      setCreateChannelError(toErrorMessage(createChannelFailure));
    } finally {
      setIsCreatingChannel(false);
    }
  };

  return (
    <main className="app">
      <h1>PulseLounge</h1>

      <div className="chat-layout">
        <ChannelPicker
          channels={channels}
          selectedChannelId={selectedChannelId}
          loading={channelsLoading}
          error={channelsError}
          newChannelName={newChannelName}
          createError={createChannelError}
          deleteError={deleteChannelError}
          isCreating={channelCreateOpen}
          isCreatingChannel={isCreatingChannel}
          deletingChannelId={deletingChannelId}
          onSelectChannel={setSelectedChannelId}
          onDeleteChannel={(channelId) => {
            void handleDeleteChannel(channelId);
          }}
          onChangeNewChannelName={setNewChannelName}
          onStartCreateChannel={() => {
            setChannelCreateOpen(true);
            setCreateChannelError(null);
          }}
          onCancelCreateChannel={() => {
            setChannelCreateOpen(false);
            setNewChannelName("");
            setCreateChannelError(null);
          }}
          onSubmitNewChannel={(event) => {
            void handleSubmitNewChannel(event);
          }}
        />

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
      </div>
    </main>
  );
}
