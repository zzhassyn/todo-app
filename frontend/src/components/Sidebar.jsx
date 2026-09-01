import { useState } from "react";
import SearchBar from "./SearchBar";

const SYSTEM_FILTERS = [
  ["all", "Все", AllIcon],
  ["today", "Сегодня", TodayIcon],
  ["active", "Активные", ActiveIcon],
  ["done", "Выполненные", DoneIcon],
  ["archive", "Архив", ArchiveIcon],
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
  user,
  theme,
  onToggleTheme,
  onLogout,
  onMoveToFolder,
  searchQuery,
  setSearchQuery,
  searchScope,
  setSearchScope,
  showScopeToggle,
}) {
  const handleDropToSystem = (key, e) => {
    e.preventDefault();
    const taskIdStr = e.dataTransfer.getData("application/x-task-id");
    if (taskIdStr && key === "all") {
      onMoveToFolder(parseInt(taskIdStr, 10), null);
    }
  };

  return (
    <aside className="sidebar">
      <div className="sidebar__brand">
        <span className="sidebar__brand-mark" aria-hidden="true" />
        Задачи
      </div>

      <div className="sidebar__search" style={{ padding: '0 var(--space-2)' }}>
        <SearchBar
          value={searchQuery}
          onChange={setSearchQuery}
          scope={searchScope}
          onScopeChange={setSearchScope}
          showScopeToggle={showScopeToggle}
        />
      </div>

      <nav className="sidebar__section" aria-label="Фильтры">
        {SYSTEM_FILTERS.map(([key, label, Icon]) => (
          <button
            key={key}
            type="button"
            className={`sidebar__item ${view.type === "system" && view.key === key ? "is-active" : ""}`}
            onClick={() => onSelectSystem(key)}
            onDragOver={(e) => {
              if (key === "all") e.preventDefault();
            }}
            onDrop={(e) => handleDropToSystem(key, e)}
          >
            <Icon />
            <span className="sidebar__item-label">{label}</span>
            {counts[key] > 0 && <span className="sidebar__item-count">{counts[key]}</span>}
          </button>
        ))}
      </nav>

      <div className="sidebar__section sidebar__section--grow">
        <div className="sidebar__heading">Мои списки</div>

        {foldersLoading ? (
          <div className="sidebar__hint">Загрузка…</div>
        ) : folders.length === 0 ? (
          <div className="sidebar__hint">Списков пока нет</div>
        ) : (
          <ul className="sidebar__folders">
            {folders.map((folder) => (
              <FolderItem
                key={folder.id}
                folder={folder}
                isActive={view.type === "folder" && view.id === folder.id}
                onSelect={() => onSelectFolder(folder)}
                onDelete={() => onDeleteFolder(folder.id)}
                onMoveToFolder={onMoveToFolder}
              />
            ))}
          </ul>
        )}

        <NewFolderControl onCreate={onCreateFolder} />
      </div>

      <div className="sidebar__footer">
        <button
          type="button"
          className="sidebar__theme-toggle"
          onClick={onToggleTheme}
          aria-label={theme === "light" ? "Включить тёмную тему" : "Включить светлую тему"}
          title={theme === "light" ? "Тёмная тема" : "Светлая тема"}
        >
          {theme === "light" ? <MoonIcon /> : <SunIcon />}
        </button>
        <span className="sidebar__user" title={user?.email}>
          {user?.email}
        </span>
        <button type="button" className="sidebar__logout" onClick={onLogout} title="Выйти">
          <LogoutIcon />
        </button>
      </div>
    </aside>
  );
}

function FolderItem({ folder, isActive, onSelect, onDelete, onMoveToFolder }) {
  const [isDragOver, setIsDragOver] = useState(false);

  return (
    <li 
      className={`sidebar__folder ${isActive ? "is-active" : ""} ${isDragOver ? "is-drag-over" : ""}`}
      onDragOver={(e) => {
        e.preventDefault();
        setIsDragOver(true);
      }}
      onDragLeave={() => setIsDragOver(false)}
      onDrop={(e) => {
        e.preventDefault();
        setIsDragOver(false);
        const taskIdStr = e.dataTransfer.getData("application/x-task-id");
        if (taskIdStr) {
          onMoveToFolder(parseInt(taskIdStr, 10), folder.id);
        }
      }}
    >
      <button type="button" className="sidebar__folder-select" onClick={onSelect}>
        <FolderIcon />
        <span className="sidebar__item-label">{folder.title}</span>
      </button>
      <button
        type="button"
        className="sidebar__folder-delete"
        onClick={(e) => {
          e.stopPropagation();
          onDelete();
        }}
        aria-label={`Удалить список «${folder.title}»`}
        title="Удалить список"
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
        <PlusIcon /> Новый список
      </button>
    );
  }

  return (
    <input
      className="sidebar__new-folder-input"
      autoFocus
      value={value}
      disabled={busy}
      placeholder="Название списка…"
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

/* ---- icons: small inline SVGs, no icon library ---- */

function AllIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <path d="M3 4.5h10M3 8h10M3 11.5h10" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
    </svg>
  );
}
function TodayIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <rect x="2" y="3" width="12" height="11" rx="1.5" fill="none" stroke="currentColor" strokeWidth="1.3" />
      <path d="M2 6.5h12" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="8" cy="10" r="1.3" fill="currentColor" />
    </svg>
  );
}
function ActiveIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <circle cx="8" cy="8" r="5.5" fill="none" stroke="currentColor" strokeWidth="1.4" />
      <path d="M8 5v3.2l2.2 1.3" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  );
}
function DoneIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <circle cx="8" cy="8" r="5.5" fill="none" stroke="currentColor" strokeWidth="1.4" />
      <path d="M5.3 8.2l1.9 1.9 3.5-4" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}
function ArchiveIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <rect x="2" y="2.5" width="12" height="3" rx="0.8" fill="none" stroke="currentColor" strokeWidth="1.3" />
      <path d="M3 5.5v7a1 1 0 001 1h8a1 1 0 001-1v-7M6.5 8.5h3" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  );
}
function FolderIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <path
        d="M2 4.2c0-.66.54-1.2 1.2-1.2h2.9l1.2 1.4h5.5c.66 0 1.2.54 1.2 1.2v6.2c0 .66-.54 1.2-1.2 1.2H3.2c-.66 0-1.2-.54-1.2-1.2V4.2z"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.3"
        strokeLinejoin="round"
      />
    </svg>
  );
}
function PlusIcon() {
  return (
    <svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true">
      <path d="M8 3v10M3 8h10" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
    </svg>
  );
}
function SunIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <circle cx="8" cy="8" r="3.2" fill="none" stroke="currentColor" strokeWidth="1.3" />
      <path
        d="M8 1.5v1.6M8 12.9v1.6M2.6 8h1.6M11.8 8h1.6M4 4l1.1 1.1M11 11l1.1 1.1M12.1 4L11 5.1M5.1 11L4 12.1"
        stroke="currentColor"
        strokeWidth="1.2"
        strokeLinecap="round"
      />
    </svg>
  );
}
function MoonIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <path
        d="M13 9.6A5.5 5.5 0 116.4 3 4.4 4.4 0 0013 9.6z"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.3"
        strokeLinejoin="round"
      />
    </svg>
  );
}
function LogoutIcon() {
  return (
    <svg viewBox="0 0 16 16" width="15" height="15" aria-hidden="true">
      <path
        d="M6.5 14H3.2a1.2 1.2 0 01-1.2-1.2V3.2A1.2 1.2 0 013.2 2h3.3M10.5 11l3-3-3-3M13.3 8H6"
        fill="none"
        stroke="currentColor"
        strokeWidth="1.3"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}
