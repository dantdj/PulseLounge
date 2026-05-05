import { cleanup, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { App } from "./App";

function createResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(body),
  } as Response;
}

const channelsResponse = [
  {
    id: 1,
    name: "general",
    created_at: "2026-03-17T17:00:00Z",
  },
  {
    id: 2,
    name: "random",
    created_at: "2026-03-17T17:05:00Z",
  },
];

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
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 1,
          author_id: 1,
          channel_id: 1,
          body: "First message",
          created_at: "2026-03-17T18:00:00Z",
          edited_at: null,
        },
      ]),
    );

    render(<App />);

    expect(await screen.findByRole("button", { name: "#general" })).toBeInTheDocument();
    expect(await screen.findByText("First message")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "#general" })).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith("/api/channels");
    expect(fetchMock).toHaveBeenCalledWith("/api/channels/1/messages");
  });

  it("shows a validation error when submitting a blank message", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse([]));

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("No messages yet.")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Send" }));

    expect(screen.getByText("Failed to send message: Message cannot be empty.")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("submits a message and clears the form", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse([]));
    fetchMock.mockResolvedValueOnce(
      createResponse({
        id: 2,
        author_id: 1,
        channel_id: 1,
        body: "A fresh note",
        created_at: "2026-03-17T18:05:00Z",
        edited_at: null,
      }, 201),
    );

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("No messages yet.")).toBeInTheDocument();
    });

    const input = screen.getByLabelText("New message");
    await user.type(input, "A fresh note");
    await user.click(screen.getByRole("button", { name: "Send" }));

    expect(await screen.findByText("A fresh note")).toBeInTheDocument();
    expect(input).toHaveValue("");

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(3);
    });

    expect(fetchMock.mock.calls[2]?.[0]).toBe("/api/channels/1/messages");
    expect(fetchMock.mock.calls[2]?.[1]).toMatchObject({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ body: "A fresh note" }),
    });
  });

  it("uploads a selected image before creating the message", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse([]));
    fetchMock.mockResolvedValueOnce(createResponse({ key: "image-key.png", url: "/media/image-key.png" }));
    fetchMock.mockResolvedValueOnce(
      createResponse({
        id: 2,
        author_id: 1,
        channel_id: 1,
        body: "With image",
        image: {
          url: "/media/image-key.png",
        },
        created_at: "2026-03-17T18:05:00Z",
        edited_at: null,
      }, 201),
    );

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("No messages yet.")).toBeInTheDocument();
    });

    const input = screen.getByLabelText("New message");
    const imageInput = screen.getByLabelText("Add image");
    const file = new File(["image-bytes"], "photo.png", { type: "image/png" });

    await user.type(input, "With image");
    await user.upload(imageInput, file);
    expect(screen.getByText("photo.png")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Send" }));

    expect(await screen.findByText("With image")).toBeInTheDocument();
    expect(screen.queryByText("photo.png")).not.toBeInTheDocument();

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(4);
    });

    const uploadOptions = fetchMock.mock.calls[2]?.[1] as RequestInit;
    expect(fetchMock.mock.calls[2]?.[0]).toBe("/api/upload");
    expect(uploadOptions).toMatchObject({
      method: "POST",
    });
    expect(uploadOptions.body).toBeInstanceOf(FormData);

    expect(fetchMock.mock.calls[3]?.[0]).toBe("/api/channels/1/messages");
    expect(fetchMock.mock.calls[3]?.[1]).toMatchObject({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ body: "With image", imageKey: "image-key.png" }),
    });
  });

  it("shows load failures from the API", async () => {
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse({ error: "database offline" }, 500));

    render(<App />);

    expect(
      await screen.findByText("Failed to load messages: database offline"),
    ).toBeInTheDocument();
  });

  it("shows channel load failures", async () => {
    fetchMock.mockResolvedValueOnce(createResponse({ error: "channels offline" }, 500));

    render(<App />);

    expect(
      await screen.findByText("Failed to load channels: channels offline"),
    ).toBeInTheDocument();
  });

  it("loads messages for the selected channel", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 1,
          author_id: 1,
          channel_id: 1,
          body: "General message",
          created_at: "2026-03-17T18:00:00Z",
          edited_at: null,
        },
      ]),
    );
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 2,
          author_id: 1,
          channel_id: 2,
          body: "Random message",
          created_at: "2026-03-17T18:10:00Z",
          edited_at: null,
        },
      ]),
    );

    render(<App />);

    expect(await screen.findByText("General message")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "#random" }));

    expect(await screen.findByText("Random message")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "#random" })).toBeInTheDocument();
    expect(fetchMock.mock.calls[2]?.[0]).toBe("/api/channels/2/messages");
  });

  it("creates a channel and selects it", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse([]));
    fetchMock.mockResolvedValueOnce(
      createResponse({
        id: 3,
        name: "ops",
        created_at: "2026-03-17T17:10:00Z",
      }, 201),
    );
    fetchMock.mockResolvedValueOnce(createResponse([]));

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("No messages yet.")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Add channel" }));

    const input = screen.getByLabelText("Channel name");
    await user.type(input, "ops");
    await user.click(screen.getByRole("button", { name: "Create channel" }));

    expect(await screen.findByRole("button", { name: "#ops" })).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "#ops" })).toBeInTheDocument();

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(4);
    });

    expect(fetchMock.mock.calls[2]?.[0]).toBe("/api/channels");
    expect(fetchMock.mock.calls[2]?.[1]).toMatchObject({
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name: "ops" }),
    });
    expect(fetchMock.mock.calls[3]?.[0]).toBe("/api/channels/3/messages");
  });

  it("deletes a selected channel and selects the next available channel", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 1,
          author_id: 1,
          channel_id: 1,
          body: "General message",
          created_at: "2026-03-17T18:00:00Z",
          edited_at: null,
        },
      ]),
    );
    fetchMock.mockResolvedValueOnce(createResponse({}, 204));
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 2,
          author_id: 1,
          channel_id: 2,
          body: "Random message",
          created_at: "2026-03-17T18:10:00Z",
          edited_at: null,
        },
      ]),
    );

    render(<App />);

    expect(await screen.findByText("General message")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Delete channel #general" }));

    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "#general" })).not.toBeInTheDocument();
    });

    expect(await screen.findByText("Random message")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "#random" })).toBeInTheDocument();
    expect(fetchMock.mock.calls[2]?.[0]).toBe("/api/channels/1");
    expect(fetchMock.mock.calls[2]?.[1]).toMatchObject({ method: "DELETE" });
    expect(fetchMock.mock.calls[3]?.[0]).toBe("/api/channels/2/messages");
  });

  it("shows delete channel failures", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse([]));
    fetchMock.mockResolvedValueOnce(createResponse({ error: "delete denied" }, 500));

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("No messages yet.")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Delete channel #general" }));

    expect(
      await screen.findByText("Failed to delete channel: delete denied"),
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "#general" })).toBeInTheDocument();
  });

  it("shows validation errors when creating a blank channel", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(createResponse([]));

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText("No messages yet.")).toBeInTheDocument();
    });

    await user.click(screen.getByRole("button", { name: "Add channel" }));
    await user.click(screen.getByRole("button", { name: "Create channel" }));

    expect(
      screen.getByText("Failed to create channel: Channel name cannot be empty."),
    ).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("edits a message inline", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 1,
          author_id: 1,
          channel_id: 1,
          body: "First message",
          created_at: "2026-03-17T18:00:00Z",
          edited_at: null,
        },
      ]),
    );
    fetchMock.mockResolvedValueOnce(createResponse({}, 204));

    render(<App />);

    expect(await screen.findByText("First message")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Edit message #1" }));

    const input = screen.getByLabelText("Edit message #1");
    await user.clear(input);
    await user.type(input, "Updated message");
    await user.click(screen.getByRole("button", { name: "Save" }));

    expect(await screen.findByText("Updated message")).toBeInTheDocument();

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(3);
    });

    expect(fetchMock.mock.calls[2]?.[0]).toBe("/api/messages/1");
    expect(fetchMock.mock.calls[2]?.[1]).toMatchObject({
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ newBody: "Updated message" }),
    });
  });

  it("shows validation errors when saving an empty edit", async () => {
    const user = userEvent.setup();
    fetchMock.mockResolvedValueOnce(createResponse(channelsResponse));
    fetchMock.mockResolvedValueOnce(
      createResponse([
        {
          id: 1,
          author_id: 1,
          channel_id: 1,
          body: "First message",
          created_at: "2026-03-17T18:00:00Z",
          edited_at: null,
        },
      ]),
    );

    render(<App />);

    expect(await screen.findByText("First message")).toBeInTheDocument();

    await user.click(screen.getByRole("button", { name: "Edit message #1" }));

    const input = screen.getByLabelText("Edit message #1");
    await user.clear(input);
    await user.click(screen.getByRole("button", { name: "Save" }));

    expect(screen.getByText("Failed to save edit: Message cannot be empty.")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });
});
