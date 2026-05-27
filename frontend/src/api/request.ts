const retryAttempts = 3;
const retryBaseDelayMs = 100;

export async function getErrorMessage(response: Response): Promise<string> {
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

export async function fetchWithRetry(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  let delayMs = retryBaseDelayMs;
  for (let attempt = 1; attempt <= retryAttempts; attempt += 1) {
    try {
      const response = init === undefined ? await fetch(input) : await fetch(input, init);
      if (response.ok || !isRetryableStatus(response.status) || attempt === retryAttempts) {
        return response;
      }
    } catch (error) {
      if (attempt === retryAttempts) {
        throw error;
      }
    }

    await delay(delayMs);
    delayMs *= 2;
  }

  throw new Error("request failed");
}

function isRetryableStatus(status: number): boolean {
  return status === 408 || status === 429 || status >= 500;
}

function delay(ms: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}
