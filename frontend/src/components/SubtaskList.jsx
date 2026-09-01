import { useState, useRef, useEffect } from "react";
import { api } from "../api/client";

export default function SubtaskList({ taskId, initialSubtasks = [], disabled = false }) {
  const [subtasks, setSubtasks] = useState(() =>
    [...initialSubtasks].sort((a, b) => a.position - b.position)
  );
  const [newTitle, setNewTitle] = useState("");
  const dragItem = useRef(null);
  const dragOverItem = useRef(null);

  // Sync if initialSubtasks changes (e.g. task switched)
  useEffect(() => {
    setSubtasks([...initialSubtasks].sort((a, b) => a.position - b.position));
  }, [initialSubtasks]);

  const handleAdd = async (e) => {
    e.preventDefault();
    const title = newTitle.trim();
    if (!title) return;

    setNewTitle("");
    const position = subtasks.length;
    
    // Optimistic UI
    const tempId = Date.now();
    const optimisticSubtask = {
      id: tempId,
      task_id: taskId,
      title,
      position,
      completed_at: null,
    };
    
    setSubtasks((prev) => [...prev, optimisticSubtask]);

    try {
      const created = await api.createSubtask(taskId, { title, position });
      setSubtasks((prev) =>
        prev.map((st) => (st.id === tempId ? created : st))
      );
    } catch (err) {
      console.error("Failed to add subtask", err);
      // Revert on failure
      setSubtasks((prev) => prev.filter((st) => st.id !== tempId));
    }
  };

  const handleToggle = async (id, isCompleted) => {
    // Optimistic UI
    setSubtasks((prev) =>
      prev.map((st) =>
        st.id === id ? { ...st, completed_at: isCompleted ? null : new Date().toISOString() } : st
      )
    );

    try {
      if (isCompleted) {
        await api.uncompleteSubtask(id);
      } else {
        await api.completeSubtask(id);
      }
    } catch (err) {
      console.error("Failed to toggle subtask", err);
      // Let's just fetch original or revert
      setSubtasks((prev) =>
        prev.map((st) =>
          st.id === id ? { ...st, completed_at: isCompleted ? new Date().toISOString() : null } : st
        )
      );
    }
  };

  const handleDelete = async (id) => {
    const backup = [...subtasks];
    setSubtasks((prev) => prev.filter((st) => st.id !== id));

    try {
      await api.deleteSubtask(id);
    } catch (err) {
      console.error("Failed to delete subtask", err);
      setSubtasks(backup);
    }
  };

  const handleSort = async () => {
    if (dragItem.current === null || dragOverItem.current === null) return;
    if (dragItem.current === dragOverItem.current) return;

    const _subtasks = [...subtasks];
    const draggedItemContent = _subtasks.splice(dragItem.current, 1)[0];
    _subtasks.splice(dragOverItem.current, 0, draggedItemContent);

    // Update positions locally
    const updatedSubtasks = _subtasks.map((st, idx) => ({ ...st, position: idx }));
    setSubtasks(updatedSubtasks);

    dragItem.current = null;
    dragOverItem.current = null;

    try {
      const ids = updatedSubtasks.map((st) => st.id);
      await api.reorderSubtasks(taskId, ids);
    } catch (err) {
      console.error("Failed to reorder subtasks", err);
    }
  };

  return (
    <div className="subtask-list">
      <div className="subtask-list__items">
        {subtasks.map((st, idx) => (
          <div
            key={st.id}
            className="subtask-row"
            draggable={!disabled}
            onDragStart={(e) => {
              dragItem.current = idx;
              e.dataTransfer.setData("application/x-subtask-idx", idx.toString());
              e.dataTransfer.effectAllowed = "move";
            }}
            onDragEnter={(e) => {
              if (e.dataTransfer.types.includes("application/x-subtask-idx")) {
                dragOverItem.current = idx;
              }
            }}
            onDragEnd={handleSort}
            onDragOver={(e) => {
              if (e.dataTransfer.types.includes("application/x-subtask-idx")) {
                e.preventDefault();
              }
            }}
          >
            <div className="subtask-row__drag-handle" aria-hidden="true">
              ⋮⋮
            </div>
            <label className="subtask-row__label">
              <input
                type="checkbox"
                className="subtask-row__checkbox"
                checked={!!st.completed_at}
                onChange={() => handleToggle(st.id, !!st.completed_at)}
                disabled={disabled}
              />
              <span className={`subtask-row__title ${st.completed_at ? "subtask-row__title--completed" : ""}`}>
                {st.title}
              </span>
            </label>
            {!disabled && (
              <button
                type="button"
                className="subtask-row__delete"
                onClick={() => handleDelete(st.id)}
                aria-label="Удалить подзадачу"
              >
                ✕
              </button>
            )}
          </div>
        ))}
      </div>
      {!disabled && (
        <form className="subtask-list__add-form" onSubmit={handleAdd}>
          <input
            type="text"
            className="subtask-list__add-input"
            value={newTitle}
            onChange={(e) => setNewTitle(e.target.value)}
            placeholder="+ Добавить подзадачу"
            maxLength={100}
          />
        </form>
      )}
    </div>
  );
}
