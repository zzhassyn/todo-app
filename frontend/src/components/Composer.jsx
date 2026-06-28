import { useMemo, useState } from "react";
import { parseTaskInput } from "../utils/parseTaskInput";

const PRIORITY_LABEL = { low: "низкий", medium: "средний", high: "высокий" };

export default function Composer({ onAdd, nextIndex, disabled, folderId, folderTitle }) {
  const [value, setValue] = useState("");
  const [dueAt, setDueAt] = useState("");
  const [showDuePicker, setShowDuePicker] = useState(false);

  const parsed = useMemo(() => parseTaskInput(value), [value]);
  const hasMarkup = parsed.tags.length > 0 || parsed.priority !== null;

  function handleSubmit(e) {
    e.preventDefault();
    if (disabled) return;
    const { title, tags, priority } = parseTaskInput(value);
    if (!title) return;

    onAdd({
      title,
      tags,
      priority: priority || undefined,
      due_at: dueAt ? new Date(dueAt).toISOString() : undefined,
      folder_id: folderId || undefined,
    });

    setValue("");
    setDueAt("");
    setShowDuePicker(false);
  }

  const placeholder = folderTitle
    ? `новая задача в «${folderTitle}»… #тег !high`
    : "новая задача… #тег !high";

  return (
    <form className="composer-wrap" onSubmit={handleSubmit}>
      <div className="composer">
        <span className="line__no line__no--ghost">{String(nextIndex).padStart(2, "0")}</span>
        <span className="line__marker line__marker--add composer__marker">+</span>
        <input
          className="composer__input"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder={placeholder}
          maxLength={140}
          disabled={disabled}
          aria-label="новая задача"
        />
        <button
          type="button"
          className={`composer__due-toggle ${dueAt ? "has-value" : ""}`}
          onClick={() => setShowDuePicker((v) => !v)}
          title="срок выполнения"
          aria-label="срок выполнения"
        >
          {dueAt ? new Date(dueAt).toLocaleDateString("ru-RU") : "⏵"}
        </button>
        <span className="composer__caret" aria-hidden="true" />
      </div>

      {showDuePicker && (
        <div className="composer__due-row">
          <input
            type="date"
            className="composer__due-input"
            value={dueAt}
            onChange={(e) => setDueAt(e.target.value)}
          />
          {dueAt && (
            <button type="button" className="link-btn" onClick={() => setDueAt("")}>
              убрать срок
            </button>
          )}
        </div>
      )}

      {hasMarkup && (
        <div className="composer__hint">
          {parsed.tags.map((t) => (
            <span key={t} className="composer__hint-tag">
              #{t}
            </span>
          ))}
          {parsed.priority && (
            <span className={`composer__hint-priority priority--${parsed.priority}`}>
              приоритет: {PRIORITY_LABEL[parsed.priority]}
            </span>
          )}
        </div>
      )}
    </form>
  );
}
