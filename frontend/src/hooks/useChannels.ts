import React from "react";
import { createChannel, deleteChannel, listChannels } from "../api/channels";
import type { Channel } from "../types";

type UseChannelsResult = {
  channels: Channel[];
  selectedChannel: Channel | null;
  selectedChannelId: number | null;
  loading: boolean;
  error: string | null;
  newChannelName: string;
  createError: string | null;
  deleteError: string | null;
  isCreating: boolean;
  isCreatingChannel: boolean;
  deletingChannelId: number | null;
  selectChannel: (channelId: number) => void;
  deleteChannelById: (channelId: number) => Promise<void>;
  setNewChannelName: (name: string) => void;
  startCreateChannel: () => void;
  cancelCreateChannel: () => void;
  submitNewChannel: (event: React.FormEvent<HTMLFormElement>) => Promise<void>;
};

function toErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "unknown error";
}

export function useChannels(): UseChannelsResult {
  const [channels, setChannels] = React.useState<Channel[]>([]);
  const [selectedChannelId, setSelectedChannelId] = React.useState<number | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [isCreatingChannel, setIsCreatingChannel] = React.useState(false);
  const [newChannelName, setNewChannelName] = React.useState("");
  const [createError, setCreateError] = React.useState<string | null>(null);
  const [deleteError, setDeleteError] = React.useState<string | null>(null);
  const [isCreating, setIsCreating] = React.useState(false);
  const [deletingChannelId, setDeletingChannelId] = React.useState<number | null>(null);

  React.useEffect(() => {
    let cancelled = false;

    async function loadChannels() {
      setLoading(true);
      setError(null);

      try {
        const result = await listChannels();
        if (cancelled) {
          return;
        }

        setChannels(result);
        setSelectedChannelId((currentChannelId) => currentChannelId ?? result[0]?.id ?? null);
      } catch (loadChannelsError: unknown) {
        if (!cancelled) {
          setError(toErrorMessage(loadChannelsError));
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    void loadChannels();

    return () => {
      cancelled = true;
    };
  }, []);

  const selectedChannel = channels.find((channel) => channel.id === selectedChannelId) ?? null;

  const startCreateChannel = React.useCallback(() => {
    setIsCreating(true);
    setCreateError(null);
  }, []);

  const cancelCreateChannel = React.useCallback(() => {
    setIsCreating(false);
    setNewChannelName("");
    setCreateError(null);
  }, []);

  const submitNewChannel = React.useCallback(async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const trimmedName = newChannelName.trim();
    if (!trimmedName) {
      setCreateError("Channel name cannot be empty.");
      return;
    }

    setIsCreatingChannel(true);
    setCreateError(null);

    try {
      const created = await createChannel(trimmedName);
      setChannels((currentChannels) => [...currentChannels, created]);
      setSelectedChannelId(created.id);
      setNewChannelName("");
      setIsCreating(false);
    } catch (createChannelFailure: unknown) {
      setCreateError(toErrorMessage(createChannelFailure));
    } finally {
      setIsCreatingChannel(false);
    }
  }, [newChannelName]);

  const deleteChannelById = React.useCallback(async (channelId: number) => {
    setDeletingChannelId(channelId);
    setDeleteError(null);

    try {
      await deleteChannel(channelId);
      setChannels((currentChannels) => {
        const nextChannels = currentChannels.filter((channel) => channel.id !== channelId);
        setSelectedChannelId((currentSelectedChannelId) =>
          currentSelectedChannelId === channelId ? nextChannels[0]?.id ?? null : currentSelectedChannelId,
        );
        return nextChannels;
      });
    } catch (deleteChannelFailure: unknown) {
      setDeleteError(toErrorMessage(deleteChannelFailure));
    } finally {
      setDeletingChannelId(null);
    }
  }, []);

  return {
    channels,
    selectedChannel,
    selectedChannelId,
    loading,
    error,
    newChannelName,
    createError,
    deleteError,
    isCreating,
    isCreatingChannel,
    deletingChannelId,
    selectChannel: setSelectedChannelId,
    deleteChannelById,
    setNewChannelName,
    startCreateChannel,
    cancelCreateChannel,
    submitNewChannel,
  };
}
