import TaskRow from "./TaskRow";
import QuickAdd from "./QuickAdd";

export default function TaskList({
  title,
  emptyHint,
  tasks,
  loading,
  isArchiveView,
  canAdd,
  contextLabel,
  allTags,
  selectedTaskId,
  onAdd,
  onToggle,
  onSelect,
  onArchive,
}) {
  return (
    <section className="tasklist">
      <header className="tasklist__header">
        <h1 className="tasklist__title">{title}</h1>
        <span className="tasklist__count">{loading ? "" : tasks.length}</span>
      </header>

      {canAdd && <QuickAdd onAdd={onAdd} allTags={allTags} contextLabel={contextLabel} />}

      <div className="tasklist__rows">
        {loading ? (
          <div className="tasklist__empty">Загрузка…</div>
        ) : tasks.length === 0 ? (
          <div className="tasklist__empty">{emptyHint}</div>
        ) : (
          tasks.map((task) => (
            <TaskRow
              key={task.id}
              task={task}
              isSelected={task.id === selectedTaskId}
              isArchiveView={isArchiveView}
              onToggle={onToggle}
              onSelect={onSelect}
              onArchive={onArchive}
            />
          ))
        )}
      </div>
    </section>
  );
}
