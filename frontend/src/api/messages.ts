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

export async function listMessages(): Promise<Message[]> {
  const response = await fetch("/api/messages");
  if (!response.ok) {
    throw new Error(await getErrorMessage(response));
  }

  return (await response.json()) as Message[];
}

export async function createMessage(body: string): Promise<Message> {
  const response = await fetch("/api/messages", {
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
