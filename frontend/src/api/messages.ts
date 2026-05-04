import type { Message } from "../types";

async function getErrorMessage(response: Response): Promise<string> {
  try {
    const data = (await response.json()) as { error?: string };
    if (data.error) {
      return data.error;
    }
  } catch {
    // Fall back to the status code when the response isn't JSON.
  }

  return `request failed with status ${response.status}`;
}

export async function listMessages(channelId: number): Promise<Message[]> {
  const response = await fetch(`/api/channels/${channelId}/messages`);
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as Message[];
}

export async function createMessage(channelId: number, body: string): Promise<Message> {
  const response = await fetch(`/api/channels/${channelId}/messages`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ body }),
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
