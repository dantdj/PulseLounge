import React from "react";
import type { Message } from "../types";

type EditableMessageItemProps = {
  message: Message;
  onEditMessage: (id: number, body: string) => Promise<string | null>;
};

export function EditableMessageItem({
  message,
  onEditMessage,
}: EditableMessageItemProps) {
  const [isEditing, setIsEditing] = React.useState(false);
  const [draftBody, setDraftBody] = React.useState(message.body);
  const [saveError, setSaveError] = React.useState<string | null>(null);
  const [isSaving, setIsSaving] = React.useState(false);

  React.useEffect(() => {
    if (!isEditing) {
      setDraftBody(message.body);
      setSaveError(null);
    }
  }, [isEditing, message.body]);

  const handleEditStart = () => {
    setDraftBody(message.body);
    setSaveError(null);
    setIsEditing(true);
  };

  const handleCancel = () => {
    setDraftBody(message.body);
    setSaveError(null);
    setIsEditing(false);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSaving(true);

    const errorMessage = await onEditMessage(message.id, draftBody);

    setIsSaving(false);

    if (errorMessage) {
      setSaveError(errorMessage);
      return;
    }

    setSaveError(null);
    setIsEditing(false);
  };

  return (
    <li className={`message-item${isEditing ? " is-editing" : ""}`}>
      {isEditing ? (
        <form className="message-edit-form" onSubmit={(event) => void handleSubmit(event)}>
          <label className="sr-only" htmlFor={`message-edit-${message.id}`}>
            Edit message #{message.id}
          </label>
          <textarea
            id={`message-edit-${message.id}`}
            value={draftBody}
            onChange={(event) => setDraftBody(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Escape") {
                event.preventDefault();
                handleCancel();
              }

              if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();
                event.currentTarget.form?.requestSubmit();
              }
            }}
            rows={3}
            disabled={isSaving}
            autoFocus
          />
          <div className="message-form-actions">
            <button type="submit" disabled={isSaving}>
              {isSaving ? "Saving..." : "Save"}
            </button>
            <button type="button" onClick={handleCancel} disabled={isSaving}>
              Cancel
            </button>
          </div>
          {saveError && <p className="error">Failed to save edit: {saveError}</p>}
        </form>
      ) : (
        <div className="message-display">
          <div className="message-copy">
            <p className="message-body">{message.body}</p>
            <p className="message-meta">
              #{message.id} · {new Date(message.created_at).toLocaleString()}
            </p>
          </div>
          <button
            type="button"
            className="message-edit-trigger"
            aria-label={`Edit message #${message.id}`}
            onClick={handleEditStart}
          >
            {/* Pencil edit icon. */}
            <svg viewBox="0 0 20 20" aria-hidden="true" focusable="false">
              <path d="M14.69 2.86a2 2 0 0 1 2.83 2.83l-8.6 8.6-3.69.86.86-3.69 8.6-8.6Zm1.41 1.42a1 1 0 0 0-1.41 0l-.8.79 1.41 1.41.8-.79a1 1 0 0 0 0-1.41ZM7.06 12.22l-.42 1.8 1.8-.42 5.64-5.64-1.41-1.41-5.61 5.67Z" />
            </svg>
          </button>
        </div>
      )}
    </li>
  );
}
