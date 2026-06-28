import { useState } from "react";

const SYSTEM_FILTERS = [
  ["all", "все"],
  ["active", "активные"],
  ["done", "выполнено"],
  ["archive", "архив"],
];

export default function Sidebar({
  view,
  onSelectSystem,
  counts,
  folders,
  foldersLoading,
  onSelectFolder,
  onCreateFolder,
  onDeleteFolder,
}) {
  return (
    <aside className="sidebar">
      <div className="sidebar__brand">
        <span className="sidebar__brand-dot" aria-hidden="true" />
        todo
      </div>

      <nav className="sidebar__section" aria-label="фильтры">
        {SYSTEM_FILTERS.map(([key, label]) => (
          <button
            key={key}
            type="button"
            className={`sidebar__item ${view.type === "system" && view.key === key ? "is-active" : ""}`}
            onClick={() => onSelectSystem(key)}
          >
            <span className="sidebar__item-icon" aria-hidden="true">
              {key === "archive" ? "○" : key === "done" ? "x" : "+"}
            </span>
            <span className="sidebar__item-label">{label}</span>
            <span className="sidebar__item-count">{counts[key]}</span>
          </button>
        ))}
      </nav>

      <div className="sidebar__divider" role="separator" />

      <div className="sidebar__section sidebar__section--grow">
        <div className="sidebar__heading">мои списки</div>

        {foldersLoading ? (
          <div className="sidebar__hint">загрузка…</div>
        ) : folders.length === 0 ? (
          <div className="sidebar__hint">списков пока нет</div>
        ) : (
          <ul className="sidebar__folders">
            {folders.map((folder) => (
              <FolderItem
                key={folder.id}
                folder={folder}
                isActive={view.type === "folder" && view.id === folder.id}
                onSelect={() => onSelectFolder(folder)}
                onDelete={() => onDeleteFolder(folder.id)}
              />
            ))}
          </ul>
        )}

        <NewFolderControl onCreate={onCreateFolder} />
      </div>
    </aside>
  );
}

function FolderItem({ folder, isActive, onSelect, onDelete }) {
  return (
    <li className={`sidebar__item sidebar__folder ${isActive ? "is-active" : ""}`}>
      <button type="button" className="sidebar__folder-select" onClick={onSelect}>
        <span className="sidebar__item-icon" aria-hidden="true">
          ▤
        </span>
        <span className="sidebar__item-label">{folder.title}</span>
      </button>
      <button
        type="button"
        className="sidebar__folder-delete"
        onClick={(e) => {
          e.stopPropagation();
          onDelete();
        }}
        aria-label={`удалить список «${folder.title}»`}
        title="удалить список"
      >
        ×
      </button>
    </li>
  );
}

function NewFolderControl({ onCreate }) {
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState("");
  const [busy, setBusy] = useState(false);

  async function commit() {
    const title = value.trim();
    if (!title) {
      setEditing(false);
      setValue("");
      return;
    }
    setBusy(true);
    try {
      await onCreate(title);
      setValue("");
      setEditing(false);
    } finally {
      setBusy(false);
    }
  }

  if (!editing) {
    return (
      <button type="button" className="sidebar__new-folder" onClick={() => setEditing(true)}>
        <span aria-hidden="true">+</span> новый список
      </button>
    );
  }

  return (
    <input
      className="sidebar__new-folder-input"
      autoFocus
      value={value}
      disabled={busy}
      placeholder="название списка…"
      maxLength={100}
      onChange={(e) => setValue(e.target.value)}
      onKeyDown={(e) => {
        if (e.key === "Enter") commit();
        if (e.key === "Escape") {
          setValue("");
          setEditing(false);
        }
      }}
      onBlur={() => {
        if (!value.trim()) setEditing(false);
      }}
    />
  );
}
