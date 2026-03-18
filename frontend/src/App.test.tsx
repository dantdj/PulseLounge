import { cleanup, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { App } from "./App";

function createResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe("App", () => {
  const fetchMock = vi.fn();

  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal("fetch", fetchMock);
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it("loads and displays messages", async () => {
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 1,
          body: "First message",
          created_at: "2026-03-17T18:00:00Z",
        },
      ]),
    );

    render(<App />);

    expect(screen.getByText("Loading messages...")).toBeInTheDocument();
    expect(await screen.findByText("First message")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith("/api/messages");
  });

  it("shows a validation error when submitting a blank message", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse([]));

    render(<App />);

    expect(await screen.findByText("No messages yet.")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Send" }));

    expect(screen.getByText("Failed to send message: Message cannot be empty.")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it("submits a message and clears the form", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse([]));
    fetchMock.mockResolvedValueOnce(
      createResponse({
        id: 2,
        body: "A fresh note",
        created_at: "2026-03-17T18:05:00Z",
      }, 201),
    );

    render(<App />);

    expect(await screen.findByText("No messages yet.")).toBeInTheDocument();

    const input = screen.getByLabelText("New message");
    await user.type(input, "A fresh note");
    await user.click(screen.getByRole("button", { name: "Send" }));

    expect(await screen.findByText("A fresh note")).toBeInTheDocument();
    expect(input).toHaveValue("");

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(2);
    });

    expect(fetchMock.mock.calls[1]?.[0]).toBe("/api/messages");
    expect(fetchMock.mock.calls[1]?.[1]).toMatchObject({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ body: "A fresh note" }),
    });
  });

  it("shows load failures from the API", async () => {
    fetchMock.mockResolvedValueOnce(createResponse({ error: "database offline" }, 500));

    render(<App />);

    expect(
      await screen.findByText("Failed to load messages: database offline"),
    ).toBeInTheDocument();
  });
});
