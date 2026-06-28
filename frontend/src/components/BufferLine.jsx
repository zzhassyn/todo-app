import { useState } from "react";

const PRIORITY_GLYPH = { low: "▾", medium: "▴", high: "▲" };

function formatDue(dueAt) {
  const d = new Date(dueAt);
  const today = new Date();
  const isOverdue = d < today;
  const label = d.toLocaleDateString("ru-RU", { day: "2-digit", month: "short" });
  return { label, isOverdue };
}

export default function BufferLine({ task, index, mode, onToggle, onArchive, onUnarchive, onEdit }) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(task.title);
  const [hovering, setHovering] = useState(false);

  function commitEdit() {
    const trimmed = draft.trim();
    setEditing(false);
    if (trimmed && trimmed !== task.title) {
      onEdit(task.id, { title: trimmed });
    } else {
      setDraft(task.title);
    }
  }

  const isArchiveView = mode === "archive";
  const due = task.due_at ? formatDue(task.due_at) : null;

  return (
    <div
      className={`line ${task.completed ? "line--done" : ""} ${isArchiveView ? "line--archived" : ""}`}
      onMouseEnter={() => setHovering(true)}
      onMouseLeave={() => setHovering(false)}
    >
      <span className="line__no">{String(index + 1).padStart(2, "0")}</span>
      <span className={`line__marker ${task.completed ? "line__marker--done" : "line__marker--add"}`}>
        {isArchiveView ? "○" : task.completed ? "−" : "+"}
      </span>

      {isArchiveView ? (
        <button
          type="button"
          className="line__checkbox line__checkbox--restore"
          aria-label="восстановить задачу"
          title="восстановить из архива"
          onClick={() => onUnarchive(task.id)}
        >
          <span className="line__checkbox-box">↺</span>
        </button>
      ) : (
        <button
          type="button"
          className="line__checkbox"
          role="checkbox"
          aria-checked={task.completed}
          aria-label={task.completed ? "снять отметку о выполнении" : "отметить выполненной"}
          onClick={() => onToggle(task)}
        >
          <span className="line__checkbox-box">{task.completed ? "x" : " "}</span>
        </button>
      )}

      <div className="line__content">
        {editing ? (
          <input
            className="line__edit-input"
            autoFocus
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={commitEdit}
            onKeyDown={(e) => {
              if (e.key === "Enter") commitEdit();
              if (e.key === "Escape") {
                setDraft(task.title);
                setEditing(false);
              }
            }}
            maxLength={100}
          />
        ) : (
          <button
            type="button"
            className="line__title"
            onClick={() => !task.completed && !isArchiveView && setEditing(true)}
            title={task.completed || isArchiveView ? undefined : "нажмите, чтобы изменить"}
          >
            <span
              className={`line__priority priority--${task.priority}`}
              title={`приоритет: ${task.priority}`}
              aria-hidden="true"
            >
              {PRIORITY_GLYPH[task.priority] ?? ""}
            </span>
            {task.title}
          </button>
        )}

        {(task.tags?.length > 0 || due) && (
          <div className="line__meta">
            {due && (
              <span className={`line__due ${due.isOverdue && !task.completed ? "line__due--overdue" : ""}`}>
                {due.label}
              </span>
            )}
            {task.tags?.map((tag) => (
              <span key={tag.id} className="line__tag">
                #{tag.name}
              </span>
            ))}
          </div>
        )}
      </div>

      <span className={`line__actions ${hovering ? "is-visible" : ""}`}>
        {!isArchiveView && (
          <button
            type="button"
            className="line__icon-btn"
            onClick={() => onArchive(task.id)}
            aria-label="архивировать задачу"
            title="архивировать"
          >
            ⌫
          </button>
        )}
      </span>
    </div>
  );
}
