import React from "react";

type MessageComposerProps = {
  value: string;
  selectedImage: File | null;
  isSubmitting: boolean;
  submitError: string | null;
  onChange: (value: string) => void;
  onImageChange: (file: File | null) => void;
  onRefresh: () => void;
  onSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
};

export function MessageComposer({
  value,
  selectedImage,
  isSubmitting,
  submitError,
  onChange,
  onImageChange,
  onRefresh,
  onSubmit,
}: MessageComposerProps) {
  const imageInputId = React.useId();
  const imageInputRef = React.useRef<HTMLInputElement>(null);

  const handleImageChange = (file: File | null) => {
    if (file === null && imageInputRef.current) {
      imageInputRef.current.value = "";
    }
    onImageChange(file);
  };

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
        <input
          id={imageInputId}
          ref={imageInputRef}
          className="sr-only"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          disabled={isSubmitting}
          onChange={(event) => handleImageChange(event.target.files?.[0] ?? null)}
        />
        <label
          className={`image-picker-button${isSubmitting ? " is-disabled" : ""}`}
          htmlFor={imageInputId}
          aria-disabled={isSubmitting}
        >
          Add image
        </label>
        <button type="submit" disabled={isSubmitting}>
          {isSubmitting ? "Sending..." : "Send"}
        </button>
        <button type="button" onClick={onRefresh} disabled={isSubmitting}>
          Refresh
        </button>
      </div>
      {selectedImage && (
        <div className="selected-image">
          <span>{selectedImage.name}</span>
          <button type="button" onClick={() => handleImageChange(null)} disabled={isSubmitting}>
            Remove
          </button>
        </div>
      )}
      {submitError && <p className="error">Failed to send message: {submitError}</p>}
    </form>
  );
}
