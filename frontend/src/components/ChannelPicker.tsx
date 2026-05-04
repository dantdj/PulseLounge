import type React from "react";
import type { Channel } from "../types";

type ChannelPickerProps = {
  channels: Channel[];
  selectedChannelId: number | null;
  loading: boolean;
  error: string | null;
  newChannelName: string;
  createError: string | null;
  deleteError: string | null;
  isCreating: boolean;
  isCreatingChannel: boolean;
  deletingChannelId: number | null;
  onSelectChannel: (channelId: number) => void;
  onDeleteChannel: (channelId: number) => void;
  onChangeNewChannelName: (name: string) => void;
  onStartCreateChannel: () => void;
  onCancelCreateChannel: () => void;
  onSubmitNewChannel: (event: React.FormEvent<HTMLFormElement>) => void;
};

export function ChannelPicker({
  channels,
  selectedChannelId,
  loading,
  error,
  newChannelName,
  createError,
  deleteError,
  isCreating,
  isCreatingChannel,
  deletingChannelId,
  onSelectChannel,
  onDeleteChannel,
  onChangeNewChannelName,
  onStartCreateChannel,
  onCancelCreateChannel,
  onSubmitNewChannel,
}: ChannelPickerProps) {
  return (
    <aside className="channel-picker" aria-label="Channels">
      <div className="channel-picker-content">
        <h2>Channels</h2>
        {loading && <p className="channel-status">Loading channels...</p>}
        {error && <p className="error">Failed to load channels: {error}</p>}
        {!loading && !error && channels.length === 0 && (
          <p className="channel-status">No channels yet.</p>
        )}
        {!loading && !error && channels.length > 0 && (
          <ul className="channel-list">
            {channels.map((channel) => (
              <li key={channel.id}>
                <div className={channel.id === selectedChannelId ? "channel-row is-selected" : "channel-row"}>
                  <button
                    type="button"
                    className="channel-select"
                    onClick={() => onSelectChannel(channel.id)}
                  >
                    <span className="channel-name">#{channel.name}</span>
                  </button>
                  <button
                    type="button"
                    className="channel-delete"
                    disabled={deletingChannelId === channel.id}
                    onClick={(event) => {
                      event.stopPropagation();
                      onDeleteChannel(channel.id);
                    }}
                    aria-label={`Delete channel #${channel.name}`}
                  >
                    {/* Trash delete icon. */}
                    <svg aria-hidden="true" focusable="false" viewBox="0 0 20 20">
                      <path d="M7 3h6l.75 1.5H17v2H3v-2h3.25L7 3Zm-2 5h10l-.7 8.05A2.2 2.2 0 0 1 12.1 18H7.9a2.2 2.2 0 0 1-2.2-1.95L5 8Zm3 1.5v6h1.5v-6H8Zm2.5 0v6H12v-6h-1.5Z" />
                    </svg>
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
        {deleteError && <p className="error">Failed to delete channel: {deleteError}</p>}
      </div>

      <div className="channel-create">
        {isCreating ? (
          <form className="channel-create-form" onSubmit={onSubmitNewChannel}>
            <label className="sr-only" htmlFor="channel-name">
              Channel name
            </label>
            <input
              id="channel-name"
              value={newChannelName}
              onChange={(event) => onChangeNewChannelName(event.target.value)}
              placeholder="Channel name"
              disabled={isCreatingChannel}
              autoFocus
            />
            <div className="channel-create-actions">
              <button type="submit" disabled={isCreatingChannel} aria-label="Create channel">
                {/* Plus create icon. */}
                <svg aria-hidden="true" focusable="false" viewBox="0 0 20 20">
                  <path d="M9 4h2v5h5v2h-5v5H9v-5H4V9h5V4Z" />
                </svg>
              </button>
              <button
                type="button"
                disabled={isCreatingChannel}
                onClick={onCancelCreateChannel}
                aria-label="Cancel channel creation"
              >
                {/* X cancel icon. */}
                <svg aria-hidden="true" focusable="false" viewBox="0 0 20 20">
                  <path d="m5.64 4.22 4.36 4.36 4.36-4.36 1.42 1.42L11.42 10l4.36 4.36-1.42 1.42L10 11.42l-4.36 4.36-1.42-1.42L8.58 10 4.22 5.64l1.42-1.42Z" />
                </svg>
              </button>
            </div>
            {createError && <p className="error">Failed to create channel: {createError}</p>}
          </form>
        ) : (
          <button
            type="button"
            className="channel-add-trigger"
            onClick={onStartCreateChannel}
            aria-label="Add channel"
          >
            {/* Plus add icon. */}
            <svg aria-hidden="true" focusable="false" viewBox="0 0 20 20">
              <path d="M9 4h2v5h5v2h-5v5H9v-5H4V9h5V4Z" />
            </svg>
          </button>
        )}
      </div>
    </aside>
  );
}
