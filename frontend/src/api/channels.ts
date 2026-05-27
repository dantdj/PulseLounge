import type { Channel } from "../types";
import { fetchWithRetry, getErrorMessage } from "./request";

export async function listChannels(): Promise<Channel[]> {
  const response = await fetchWithRetry("/api/channels");
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as Channel[];
}

export async function createChannel(name: string): Promise<Channel> {
  const response = await fetch("/api/channels", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ name }),
  });
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as Channel;
}

export async function deleteChannel(id: number): Promise<void> {
  const response = await fetch(`/api/channels/${id}`, {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }
}
