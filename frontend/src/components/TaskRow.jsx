import { useState } from "react";
import { formatShortDate, isOverdue } from "../utils/date";
import PriorityDot from "./PriorityDot";
import TagBadge from "./TagBadge";

export default function TaskRow({ task, isSelected, isArchiveView, onToggle, onSelect, onArchive, onReorder }) {
  const [dragOverPos, setDragOverPos] = useState(null);
  
  const due = task.due_at ? new Date(task.due_at) : null;
  const overdue = due && !task.completed && isOverdue(due);

  return (
    <div
      draggable
      className={`task-row ${isSelected ? "is-selected" : ""} ${task.completed ? "is-completed" : ""} ${dragOverPos ? `is-drag-over-${dragOverPos}` : ""}`}
      onClick={() => onSelect(task)}
      onDragStart={(e) => {
        e.dataTransfer.setData("application/x-task-id", task.id.toString());
      }}
      onDragOver={(e) => {
        e.preventDefault();
        const rect = e.currentTarget.getBoundingClientRect();
        const y = e.clientY - rect.top;
        setDragOverPos(y < rect.height / 2 ? "before" : "after");
      }}
      onDragLeave={() => setDragOverPos(null)}
      onDrop={(e) => {
        e.preventDefault();
        setDragOverPos(null);
        const draggedIdStr = e.dataTransfer.getData("application/x-task-id");
        if (draggedIdStr && onReorder) {
          const draggedId = parseInt(draggedIdStr, 10);
          if (draggedId !== task.id) {
            onReorder(draggedId, task.id, dragOverPos);
          }
        }
      }}
    >
      <button
        type="button"
        className="task-row__check"
        role="checkbox"
        aria-checked={task.completed}
        aria-label={task.completed ? "Снять отметку о выполнении" : "Отметить выполненной"}
        onClick={(e) => {
          e.stopPropagation();
          onToggle(task);
        }}
      >
        {task.completed && (
          <svg viewBox="0 0 16 16" width="10" height="10" aria-hidden="true">
            <path
              d="M3 8.5L6.2 11.5L13 4"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        )}
      </button>

      <div className="task-row__body">
        <div className="task-row__title-line">
          <PriorityDot priority={task.priority} />
          <span className="task-row__title">{task.title}</span>
        </div>

        {(task.tags?.length > 0 || due) && (
          <div className="task-row__meta">
            {due && (
              <span className={`task-row__due ${overdue ? "is-overdue" : ""}`}>
                {formatShortDate(due)}
              </span>
            )}
            {task.tags?.map((tag) => (
              <TagBadge key={tag.id} name={tag.name} />
            ))}
          </div>
        )}
      </div>

      {!isArchiveView && (
        <button
          type="button"
          className="task-row__archive"
          onClick={(e) => {
            e.stopPropagation();
            onArchive(task.id);
          }}
          aria-label="Архивировать задачу"
          title="Архивировать"
        >
          <svg viewBox="0 0 16 16" width="14" height="14" aria-hidden="true">
            <path
              d="M2 4h12M3 4v9a1 1 0 001 1h8a1 1 0 001-1V4M6 7h4"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.4"
              strokeLinecap="round"
            />
            <path d="M5 4l1-2h4l1 2" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </button>
      )}
    </div>
  );
}
