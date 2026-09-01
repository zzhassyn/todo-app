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
  batchSelectedIds,
  onToggleBatchSelect,
  onBulkComplete,
  onBulkArchive,
  onBulkMove,
  folders, // needed for bulk move
  onAdd,
  onToggle,
  onSelect,
  onArchive,
  onReorder,
}) {
  const isSelectionMode = batchSelectedIds?.size > 0;

  return (
    <section className="tasklist">
      <header className="tasklist__header">
        <h1 className="tasklist__title">{title}</h1>
        <span className="tasklist__count">{loading ? "" : tasks.length}</span>
      </header>

      {isSelectionMode && (
        <div className="tasklist__bulk-actions" style={{
          padding: '12px', background: 'var(--color-bg-elevated)', border: '1px solid var(--color-border)',
          borderRadius: '8px', marginBottom: '16px', display: 'flex', gap: '8px', alignItems: 'center'
        }}>
          <span>Выбрано: {batchSelectedIds.size}</span>
          <button className="btn-secondary" onClick={onBulkComplete}>Выполнить</button>
          <button className="btn-secondary" onClick={onBulkArchive}>В архив</button>
          
          <select className="btn-secondary" onChange={(e) => {
            if (e.target.value) onBulkMove(e.target.value);
            e.target.value = "";
          }}>
            <option value="">Переместить в...</option>
            {folders?.map(f => (
              <option key={f.id} value={f.id}>{f.title}</option>
            ))}
          </select>

          <div style={{flex: 1}}></div>
          <button className="btn-secondary" onClick={() => onToggleBatchSelect("CLEAR")}>Отмена</button>
        </div>
      )}

      {!isSelectionMode && canAdd && <QuickAdd onAdd={onAdd} allTags={allTags} contextLabel={contextLabel} />}

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
              isBatchSelected={batchSelectedIds?.has(task.id)}
              isArchiveView={isArchiveView}
              onToggle={onToggle}
              onSelect={onSelect}
              onArchive={onArchive}
              onReorder={onReorder}
              onToggleBatchSelect={onToggleBatchSelect}
            />
          ))
        )}
      </div>
    </section>
  );
}
