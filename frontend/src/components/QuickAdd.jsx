import { useRef, useState } from "react";
import DatePicker from "./DatePicker";
import PriorityPicker from "./PriorityPicker";
import TagPicker from "./TagPicker";
import PriorityDot from "./PriorityDot";
import TagBadge from "./TagBadge";
import { formatShortDate } from "../utils/date";

export default function QuickAdd({ onAdd, allTags, contextLabel }) {
  const [title, setTitle] = useState("");
  const [dueAt, setDueAt] = useState(null);
  const [priority, setPriority] = useState(null);
  const [tags, setTags] = useState([]);
  const [openPicker, setOpenPicker] = useState(null); // "date" | "priority" | "tags" | null

  const dateBtnRef = useRef(null);
  const priorityBtnRef = useRef(null);
  const tagsBtnRef = useRef(null);

  const hasStarted = title.trim().length > 0;

  function reset() {
    setTitle("");
    setDueAt(null);
    setPriority(null);
    setTags([]);
    setOpenPicker(null);
  }

  function handleSubmit(e) {
    e.preventDefault();
    const trimmed = title.trim();
    if (!trimmed) return;

    onAdd({
      title: trimmed,
      due_at: dueAt ? dueAt.toISOString() : undefined,
      priority: priority || undefined,
      tags: tags.length > 0 ? tags : undefined,
    });
    reset();
  }

  function toggleTag(name) {
    setTags((prev) => (prev.includes(name) ? prev.filter((t) => t !== name) : [...prev, name]));
  }

  return (
    <form className="quickadd" onSubmit={handleSubmit}>
      <div className="quickadd__row">
        <span className="quickadd__bullet" aria-hidden="true" />
        <input
          className="quickadd__input"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder={contextLabel ? `Добавить задачу в «${contextLabel}»…` : "Добавить новую задачу…"}
          maxLength={140}
        />
        <button type="submit" className="quickadd__submit" disabled={!hasStarted}>
          Добавить
        </button>
      </div>

      {hasStarted && (
        <div className="quickadd__triggers">
          <div className="quickadd__trigger-wrap">
            <button
              type="button"
              ref={dateBtnRef}
              className={`quickadd__trigger ${dueAt ? "is-set" : ""}`}
              onClick={() => setOpenPicker(openPicker === "date" ? null : "date")}
            >
              <CalendarIcon />
              {dueAt ? formatShortDate(dueAt) : "Дата"}
            </button>
            {openPicker === "date" && (
              <DatePicker
                value={dueAt}
                anchorRef={dateBtnRef}
                onClose={() => setOpenPicker(null)}
                onChange={(d) => {
                  setDueAt(d);
                  setOpenPicker(null);
                }}
              />
            )}
          </div>

          <div className="quickadd__trigger-wrap">
            <button
              type="button"
              ref={priorityBtnRef}
              className={`quickadd__trigger ${priority ? "is-set" : ""}`}
              onClick={() => setOpenPicker(openPicker === "priority" ? null : "priority")}
            >
              {priority ? <PriorityDot priority={priority} /> : <FlagIcon />}
              {priority ? { low: "Низкий", medium: "Средний", high: "Высокий" }[priority] : "Приоритет"}
            </button>
            {openPicker === "priority" && (
              <PriorityPicker
                value={priority}
                anchorRef={priorityBtnRef}
                onClose={() => setOpenPicker(null)}
                onChange={setPriority}
              />
            )}
          </div>

          <div className="quickadd__trigger-wrap">
            <button
              type="button"
              ref={tagsBtnRef}
              className={`quickadd__trigger ${tags.length > 0 ? "is-set" : ""}`}
              onClick={() => setOpenPicker(openPicker === "tags" ? null : "tags")}
            >
              <TagIcon />
              Теги
            </button>
            {openPicker === "tags" && (
              <TagPicker
                allTags={allTags}
                selected={tags}
                anchorRef={tagsBtnRef}
                onClose={() => setOpenPicker(null)}
                onToggle={toggleTag}
              />
            )}
          </div>

          {tags.length > 0 && (
            <div className="quickadd__tags">
              {tags.map((name) => (
                <TagBadge key={name} name={name} onRemove={() => toggleTag(name)} />
              ))}
            </div>
          )}
        </div>
      )}
    </form>
  );
}

function CalendarIcon() {
  return (
    <svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
      <rect x="2" y="3" width="12" height="11" rx="1.5" fill="none" stroke="currentColor" strokeWidth="1.3" />
      <path d="M2 6.5h12M5 1.5v3M11 1.5v3" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  );
}

function FlagIcon() {
  return (
    <svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
      <path
        d="M4 1.5v13M4 2h7l-2 2.5L11 7H4"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.3"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function TagIcon() {
  return (
    <svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
      <path
        d="M2 2h5.5L14 8.5 8.5 14 2 7.5V2z"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.3"
        strokeLinejoin="round"
      />
      <circle cx="4.8" cy="4.8" r="0.9" fill="currentColor" />
    </svg>
  );
}
