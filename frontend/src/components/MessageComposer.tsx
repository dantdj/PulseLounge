import React from "react";

type MessageComposerProps = {
  value: string;
  isSubmitting: boolean;
  submitError: string | null;
  onChange: (value: string) => void;
  onRefresh: () => void;
  onSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
};

export function MessageComposer({
  value,
  isSubmitting,
  submitError,
  onChange,
  onRefresh,
  onSubmit,
}: MessageComposerProps) {
  return (
    <form className="message-form" onSubmit={onSubmit}>
      <label htmlFor="message-body">New message</label>
      <textarea
        id="message-body"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            event.currentTarget.form?.requestSubmit();
          }
        }}
        rows={3}
        placeholder="Type a message..."
        disabled={isSubmitting}
      />
      <div className="message-form-actions">
        <button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Sending..." : "Send"}
        </button>
        <button type="button" onClick={onRefresh} disabled={isSubmitting}>
          Refresh
        </button>
      </div>
      {submitError && <p className="error">Failed to send message: {submitError}</p>}
    </form>
  );
}
