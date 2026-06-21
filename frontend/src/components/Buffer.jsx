import BufferLine from "./BufferLine";
import Composer from "./Composer";

export default function Buffer({ tasks, loading, filter, onAdd, onToggle, onDelete, onEdit }) {
  return (
    <div className="buffer">
      <div className="buffer__titlebar">
        <span className="buffer__filename">buffer.todo</span>
        <span className="buffer__modified" aria-hidden={tasks.length === 0}>
          {tasks.length > 0 ? "[+]" : ""}
        </span>
      </div>

      <div className="buffer__body">
        {loading ? (
          <div className="buffer__empty">читаю буфер…</div>
        ) : tasks.length === 0 ? (
          <div className="buffer__empty">
            {filter === "done"
              ? "ничего не выполнено. пока."
              : filter === "active"
              ? "активных задач нет — всё сделано."
              : "буфер пуст. начните печатать ниже."}
          </div>
        ) : (
          tasks.map((task, i) => (
            <BufferLine
              key={task.id}
              task={task}
              index={i}
              onToggle={onToggle}
              onDelete={onDelete}
              onEdit={onEdit}
            />
          ))
        )}

        {filter === "all" && (
          <Composer onAdd={onAdd} nextIndex={tasks.length + 1} disabled={loading} />
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
