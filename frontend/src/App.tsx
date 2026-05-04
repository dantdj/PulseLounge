import { ChannelPicker } from "./components/ChannelPicker";
import { ChatView } from "./components/ChatView";
import { useChannels } from "./hooks/useChannels";

export function App() {
  const channels = useChannels();

  return (
    <main className="app">
      <h1>PulseLounge</h1>

      <div className="chat-layout">
        <ChannelPicker
          channels={channels.channels}
          selectedChannelId={channels.selectedChannelId}
          loading={channels.loading}
          error={channels.error}
          newChannelName={channels.newChannelName}
          createError={channels.createError}
          deleteError={channels.deleteError}
          isCreating={channels.isCreating}
          isCreatingChannel={channels.isCreatingChannel}
          deletingChannelId={channels.deletingChannelId}
          onSelectChannel={channels.selectChannel}
          onDeleteChannel={(channelId) => {
            void channels.deleteChannelById(channelId);
          }}
          onChangeNewChannelName={channels.setNewChannelName}
          onStartCreateChannel={channels.startCreateChannel}
          onCancelCreateChannel={channels.cancelCreateChannel}
          onSubmitNewChannel={(event) => {
            void channels.submitNewChannel(event);
          }}
        />

        <ChatView
          selectedChannel={channels.selectedChannel}
          selectedChannelId={channels.selectedChannelId}
        />
      </div>
    </main>
  );
}
