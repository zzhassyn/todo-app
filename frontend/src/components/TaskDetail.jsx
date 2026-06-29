import { useRef, useState } from "react";
import DatePicker from "./DatePicker";
import PriorityPicker from "./PriorityPicker";
import TagPicker from "./TagPicker";
import PriorityDot from "./PriorityDot";
import TagBadge from "./TagBadge";
import SubtaskList from "./SubtaskList";

const PRIORITY_LABEL = { low: "Низкий", medium: "Средний", high: "Высокий" };

function formatFull(dateStr) {
  return new Date(dateStr).toLocaleString("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

// The caller renders this with key={task.id} (see App.jsx), so React
// remounts it — and re-initializes title/description from the new task —
// whenever a different task is selected. That's simpler and more correct
// than syncing local state from props in an effect: it also means typing
// into the title/description fields is never clobbered by an unrelated
// background refresh of the *same* task (only switching tasks resets the
// draft).
export default function TaskDetail({ task, allTags, onClose, onPatch, onArchive, onUnarchive, onPermanentlyDelete }) {
  const [title, setTitle] = useState(task.title);
  const [description, setDescription] = useState(task.description || "");
  const [openPicker, setOpenPicker] = useState(null);

  const dateBtnRef = useRef(null);
  const priorityBtnRef = useRef(null);
  const tagsBtnRef = useRef(null);

  function commitTitle() {
    const trimmed = title.trim();
    if (trimmed && trimmed !== task.title) {
      onPatch(task.id, { title: trimmed });
    } else {
      setTitle(task.title);
    }
  }

  function commitDescription() {
    const trimmed = description.trim();
    if (trimmed !== (task.description || "")) {
      onPatch(task.id, { description: trimmed || null });
    }
  }

  const tagNames = task.tags?.map((t) => t.name) ?? [];

  function toggleTag(name) {
    const next = tagNames.includes(name)
      ? tagNames.filter((t) => t !== name)
      : [...tagNames, name];
    onPatch(task.id, { tags: next });
  }

  const due = task.due_at ? new Date(task.due_at) : null;
  const isArchived = Boolean(task.archived_at);

  return (
    <aside className="taskdetail">
      <div className="taskdetail__topbar">
        <button type="button" className="taskdetail__close" onClick={onClose} aria-label="Закрыть панель">
          <svg viewBox="0 0 16 16" width="14" height="14" aria-hidden="true">
            <path
              d="M3 3l10 10M13 3L3 13"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.4"
              strokeLinecap="round"
            />
          </svg>
        </button>
      </div>

      <div className="taskdetail__body">
        <textarea
          className="taskdetail__title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          onBlur={commitTitle}
          rows={1}
          maxLength={140}
          disabled={isArchived}
        />

        <textarea
          className="taskdetail__description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          onBlur={commitDescription}
          placeholder="Добавить описание…"
          rows={4}
          maxLength={1000}
          disabled={isArchived}
        />

        <div className="taskdetail__field taskdetail__field--subtasks">
          <span className="taskdetail__field-label">Чек-лист</span>
          <div className="taskdetail__field-control-wrap" style={{ display: 'block' }}>
            <SubtaskList
              taskId={task.id}
              initialSubtasks={task.subtasks || []}
              disabled={isArchived}
            />
          </div>
        </div>

        <div className="taskdetail__field">
          <span className="taskdetail__field-label">Приоритет</span>
          <div className="taskdetail__field-control-wrap">
            <button
              type="button"
              ref={priorityBtnRef}
              className="taskdetail__field-control"
              onClick={() => setOpenPicker(openPicker === "priority" ? null : "priority")}
              disabled={isArchived}
            >
              <PriorityDot priority={task.priority} />
              {PRIORITY_LABEL[task.priority] ?? "Не задан"}
            </button>
            {openPicker === "priority" && (
              <PriorityPicker
                value={task.priority}
                anchorRef={priorityBtnRef}
                onClose={() => setOpenPicker(null)}
                onChange={(p) => onPatch(task.id, { priority: p })}
              />
            )}
          </div>
        </div>

        <div className="taskdetail__field">
          <span className="taskdetail__field-label">Срок</span>
          <div className="taskdetail__field-control-wrap">
            <button
              type="button"
              ref={dateBtnRef}
              className="taskdetail__field-control"
              onClick={() => setOpenPicker(openPicker === "date" ? null : "date")}
              disabled={isArchived}
            >
              {due ? due.toLocaleDateString("ru-RU", { day: "numeric", month: "long" }) : "Не задан"}
            </button>
            {openPicker === "date" && (
              <DatePicker
                value={due}
                anchorRef={dateBtnRef}
                onClose={() => setOpenPicker(null)}
                onChange={(d) => onPatch(task.id, { due_at: d ? d.toISOString() : null })}
              />
            )}
          </div>
        </div>

        <div className="taskdetail__field taskdetail__field--tags">
          <span className="taskdetail__field-label">Теги</span>
          <div className="taskdetail__field-control-wrap">
            <div className="taskdetail__tags">
              {tagNames.map((name) => (
                <TagBadge key={name} name={name} onRemove={isArchived ? undefined : () => toggleTag(name)} />
              ))}
              {!isArchived && (
                <button
                  type="button"
                  ref={tagsBtnRef}
                  className="taskdetail__add-tag"
                  onClick={() => setOpenPicker(openPicker === "tags" ? null : "tags")}
                >
                  + Тег
                </button>
              )}
            </div>
            {openPicker === "tags" && (
              <TagPicker
                allTags={allTags}
                selected={tagNames}
                anchorRef={tagsBtnRef}
                onClose={() => setOpenPicker(null)}
                onToggle={toggleTag}
              />
            )}
          </div>
        </div>

        <div className="taskdetail__history">
          <span className="taskdetail__field-label">История</span>
          <ul className="taskdetail__history-list">
            <li>Создана {formatFull(task.created_at)}</li>
            {task.completed_at && <li>Выполнена {formatFull(task.completed_at)}</li>}
            {task.archived_at && <li>Архивирована {formatFull(task.archived_at)}</li>}
          </ul>
        </div>
      </div>

      <div className="taskdetail__footer">
        {isArchived ? (
          <>
            <button type="button" className="btn btn--secondary" onClick={() => onUnarchive(task.id)}>
              Восстановить
            </button>
            <button
              type="button"
              className="btn btn--danger-ghost"
              onClick={() => {
                if (window.confirm("Удалить задачу навсегда? Это действие нельзя отменить.")) {
                  onPermanentlyDelete(task.id);
                }
              }}
            >
              Удалить навсегда
            </button>
          </>
        ) : (
          <button type="button" className="btn btn--secondary" onClick={() => onArchive(task.id)}>
            Архивировать
          </button>
        )}
      </div>
    </aside>
  );
}
