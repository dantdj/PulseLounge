import React from "react";
import { createMessage, editMessage, listMessages, uploadImage } from "../api/messages";
import type { Message } from "../types";

type UseMessagesResult = {
  messages: Message[];
  loading: boolean;
  error: string | null;
  submitError: string | null;
  isSubmitting: boolean;
  loadMessages: () => Promise<void>;
  submitMessage: (body: string, imageFile: File | null) => Promise<boolean>;
  saveEditedMessage: (id: number, body: string) => Promise<string | null>;
};

function toErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "unknown error";
}

export function useMessages(channelId: number | null): UseMessagesResult {
  const [messages, setMessages] = React.useState<Message[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [loadedChannelId, setLoadedChannelId] = React.useState<number | null>(null);
  const [error, setError] = React.useState<string | null>(null);
  const [submitError, setSubmitError] = React.useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  const loadMessages = React.useCallback(async () => {
    if (channelId === null) {
      setMessages([]);
      setLoading(false);
      setError(null);
      setLoadedChannelId(null);
      return;
    }

    setLoading(true);
    setError(null);
    setMessages([]);

    try {
      const result = await listMessages(channelId);
      setMessages(result);
      setLoadedChannelId(channelId);
    } catch (loadError: unknown) {
      setError(toErrorMessage(loadError));
      setLoadedChannelId(channelId);
    } finally {
      setLoading(false);
    }
  }, [channelId]);

  React.useEffect(() => {
    void loadMessages();
  }, [loadMessages]);

  const submitMessage = React.useCallback(async (body: string, imageFile: File | null) => {
    if (channelId === null) {
      setSubmitError("Choose a channel before sending a message.");
      return false;
    }

    const trimmedBody = body.trim();
    if (!trimmedBody) {
      setSubmitError("Message cannot be empty.");
      return false;
    }

    setSubmitError(null);
    setIsSubmitting(true);

    try {
      const uploadedImage = imageFile ? await uploadImage(imageFile) : null;
      const created = await createMessage(channelId, trimmedBody, uploadedImage?.key);
      setMessages((currentMessages) => [...currentMessages, created]);
      return true;
    } catch (submitMessageError: unknown) {
      setSubmitError(toErrorMessage(submitMessageError));
      return false;
    } finally {
      setIsSubmitting(false);
    }
  }, [channelId]);

  const saveEditedMessage = React.useCallback(async (id: number, body: string) => {
    const trimmedBody = body.trim();
    if (!trimmedBody) {
      return "Message cannot be empty.";
    }

    try {
      await editMessage(id, trimmedBody);
      setMessages((currentMessages) =>
        currentMessages.map((message) =>
          message.id === id ? { ...message, body: trimmedBody } : message,
        ),
      );
      return null;
    } catch (saveEditedMessageError: unknown) {
      return toErrorMessage(saveEditedMessageError);
    }
  }, []);

  return {
    messages,
    loading: loading || (channelId !== null && loadedChannelId !== channelId),
    error,
    submitError,
    isSubmitting,
    loadMessages,
    submitMessage,
    saveEditedMessage,
  };
}
