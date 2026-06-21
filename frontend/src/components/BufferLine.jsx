import { useState } from "react";

export default function BufferLine({ task, index, onToggle, onDelete, onEdit }) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(task.title);
  const [hovering, setHovering] = useState(false);

  function commitEdit() {
    const trimmed = draft.trim();
    setEditing(false);
    if (trimmed && trimmed !== task.title) {
      onEdit(task.id, trimmed);
    } else {
      setDraft(task.title);
    }
  }

  return (
    <div
      className={`line ${task.completed ? "line--done" : ""}`}
      onMouseEnter={() => setHovering(true)}
      onMouseLeave={() => setHovering(false)}
    >
      <span className="line__no">{String(index + 1).padStart(2, "0")}</span>
      <span className={`line__marker ${task.completed ? "line__marker--done" : "line__marker--add"}`}>
        {task.completed ? "−" : "+"}
      </span>

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
          onClick={() => !task.completed && setEditing(true)}
          title={task.completed ? undefined : "нажмите, чтобы изменить"}
        >
          {task.title}
        </button>
      )}

      <span className={`line__actions ${hovering ? "is-visible" : ""}`}>
        <button
          type="button"
          className="line__icon-btn"
          onClick={() => onDelete(task.id)}
          aria-label="удалить задачу"
          title="удалить"
        >
          ×
        </button>
      </span>
    </div>
  );
}
