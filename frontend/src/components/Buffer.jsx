import BufferLine from "./BufferLine";
import Composer from "./Composer";

export default function Buffer({
  view,
  title,
  emptyHint,
  tasks,
  loading,
  onAdd,
  onToggle,
  onArchive,
  onUnarchive,
  onEdit,
}) {
  const isArchiveView = view.type === "system" && view.key === "archive";

  return (
    <div className="buffer">
      <div className="buffer__titlebar">
        <span className="buffer__filename">{title}</span>
        <span className="buffer__modified" aria-hidden={tasks.length === 0}>
          {tasks.length > 0 ? "[+]" : ""}
        </span>
      </div>

      <div className="buffer__body">
        {loading ? (
          <div className="buffer__empty">читаю…</div>
        ) : tasks.length === 0 ? (
          <div className="buffer__empty">{emptyHint}</div>
        ) : (
          tasks.map((task, i) => (
            <BufferLine
              key={task.id}
              task={task}
              index={i}
              mode={isArchiveView ? "archive" : "normal"}
              onToggle={onToggle}
              onArchive={onArchive}
              onUnarchive={onUnarchive}
              onEdit={onEdit}
            />
          ))
        )}

        {!isArchiveView && (
          <Composer
            onAdd={onAdd}
            nextIndex={tasks.length + 1}
            disabled={loading}
            folderId={view.type === "folder" ? view.id : undefined}
            folderTitle={view.type === "folder" ? view.title : undefined}
          />
        )}

        <EndOfBuffer />
      </div>
    </div>
  );
}

function EndOfBuffer() {
  const rows = Array.from({ length: 14 });
  return (
    <div className="buffer__eof" aria-hidden="true">
      {rows.map((_, i) => (
        <div className="eof-line" key={i}>
          <span className="eof-line__no">~</span>
        </div>
      ))}
    </div>
  );
}
