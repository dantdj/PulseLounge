import type { Message } from "../types";
import { fetchWithRetry, getErrorMessage } from "./request";

export async function listMessages(channelId: number): Promise<Message[]> {
  const response = await fetchWithRetry(`/api/channels/${channelId}/messages`);
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as Message[];
}

type UploadImageResponse = {
  key: string;
  url: string;
};

export async function uploadImage(file: File): Promise<UploadImageResponse> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch("/api/upload", {
    method: "POST",
    body: formData,
  });
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as UploadImageResponse;
}

export async function createMessage(
  channelId: number,
  body: string,
  imageKey?: string,
): Promise<Message> {
  const response = await fetch(`/api/channels/${channelId}/messages`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(imageKey ? { body, imageKey } : { body }),
  });
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as Message;
}

export async function editMessage(id: number, newBody: string): Promise<void> {
  const response = await fetch(`/api/messages/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ newBody }),
  });
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }
}
