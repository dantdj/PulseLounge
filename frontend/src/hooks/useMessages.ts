import React from "react";
import { createMessage, listMessages } from "../api/messages";
import type { Message } from "../types";

type UseMessagesResult = {
  messages: Message[];
  loading: boolean;
  error: string | null;
  submitError: string | null;
  isSubmitting: boolean;
  loadMessages: () => Promise<void>;
  submitMessage: (body: string) => Promise<boolean>;
};

function toErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "unknown error";
}

export function useMessages(): UseMessagesResult {
  const [messages, setMessages] = React.useState<Message[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [submitError, setSubmitError] = React.useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = React.useState(false);

  const loadMessages = React.useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const result = await listMessages();
      setMessages(result);
    } catch (loadError: unknown) {
      setError(toErrorMessage(loadError));
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    void loadMessages();
  }, [loadMessages]);

  const submitMessage = React.useCallback(async (body: string) => {
    const trimmedBody = body.trim();
    if (!trimmedBody) {
      setSubmitError("Message cannot be empty.");
      return false;
    }

    setSubmitError(null);
    setIsSubmitting(true);

    try {
      const created = await createMessage(trimmedBody);
      setMessages((currentMessages) => [...currentMessages, created]);
      return true;
    } catch (submitMessageError: unknown) {
      setSubmitError(toErrorMessage(submitMessageError));
      return false;
    } finally {
      setIsSubmitting(false);
    }
  }, []);

  return {
    messages,
    loading,
    error,
    submitError,
    isSubmitting,
    loadMessages,
    submitMessage,
  };
}
