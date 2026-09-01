import { useCallback, useMemo, useState } from "react";
import { AuthProvider, useAuth } from "./context/AuthContext";
import { ApiError } from "./api/client";
import { isToday } from "./utils/date";
import { useTheme } from "./utils/useTheme";
import { useTasks } from "./hooks/useTasks";
import { useFolders } from "./hooks/useFolders";
import { useTaskSearch } from "./hooks/useTaskSearch";
import { useKeyboardShortcuts } from "./hooks/useKeyboardShortcuts";
import AuthScreen from "./components/AuthScreen";
import Sidebar from "./components/Sidebar";
import TaskList from "./components/TaskList";
import TaskDetail from "./components/TaskDetail";
import Toast from "./components/Toast";
import SearchBar from "./components/SearchBar";
import "./app.css";

const SYSTEM_LABELS = {
  all: "Все задачи",
  today: "Сегодня",
  active: "Активные",
  done: "Выполненные",
  archive: "Архив",
};

const SYSTEM_EMPTY_HINTS = {
  all: "Задач пока нет. Добавьте первую ниже.",
  today: "На сегодня ничего не запланировано.",
  active: "Активных задач нет — всё сделано.",
  done: "Пока ничего не выполнено.",
  archive: "Архив пуст.",
};

function TodoApp() {
  const { user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();

  // `view` is the single source of truth for "what is currently shown":
  // either one of the system filters, or a specific folder. Keeping it as
  // one discriminated value (rather than separate `filter` + `folderId`
  // state) means the sidebar's highlighted item and the list's contents
  // can never disagree with each other.
  const [view, setView] = useState({ type: "system", key: "all" });
  const [selectedTaskId, setSelectedTaskId] = useState(null);
  const [error, setError] = useState(null);

  const reportError = useCallback((err, fallback) => {
    setError(err instanceof ApiError ? err.message : fallback);
  }, []);

  const {
    optTasks,
    optFolderTasks,
    optArchivedTasks,
    loading,
    tags,
    batchSelectedIds,
    setBatchSelectedIds,
    handleAdd,
    handleToggle,
    handleArchive,
    handleUnarchive,
    handlePermanentlyDelete,
    handleEdit,
    handleMoveToFolder,
    handleReorderTask,
    handleBulkComplete,
    handleBulkArchive,
    handleBulkMove,
    handleToggleBatchSelect,
  } = useTasks({ view, reportError });

  const { folders, foldersLoading, handleCreateFolder, handleDeleteFolder: rawDeleteFolder } =
    useFolders({ reportError });

  // Wrap folder deletion to also reset view if the deleted folder is open
  const handleDeleteFolder = useCallback(
    async (id) => {
      await rawDeleteFolder(id);
      setView((prev) => (prev.type === "folder" && prev.id === id ? { type: "system", key: "all" } : prev));
    },
    [rawDeleteFolder]
  );

  const { searchQuery, setSearchQuery, searchScope, setSearchScope, searchResults, isSearching } =
    useTaskSearch({ view, reportError });

  const handleSelectSystem = useCallback((key) => {
    setView({ type: "system", key });
    setSelectedTaskId(null);
    setBatchSelectedIds(new Set());
    setSearchQuery("");
  }, [setBatchSelectedIds, setSearchQuery]);

  const handleSelectFolder = useCallback((folder) => {
    setView({ type: "folder", id: folder.id, title: folder.title });
    setSelectedTaskId(null);
    setBatchSelectedIds(new Set());
    setSearchQuery("");
  }, [setBatchSelectedIds, setSearchQuery]);

  const counts = useMemo(() => {
    const done = optTasks.filter((t) => t.completed).length;
    const today = optTasks.filter((t) => t.due_at && isToday(new Date(t.due_at))).length;
    return {
      all: optTasks.length,
      today,
      done,
      active: optTasks.length - done,
      archive: optArchivedTasks.length,
    };
  }, [optTasks, optArchivedTasks]);

  const visibleTasks = useMemo(() => {
    if (view.type === "folder") return optFolderTasks;
    if (view.key === "archive") return optArchivedTasks;
    if (view.key === "today") return optTasks.filter((t) => t.due_at && isToday(new Date(t.due_at)));
    if (view.key === "active") return optTasks.filter((t) => !t.completed);
    if (view.key === "done") return optTasks.filter((t) => t.completed);
    return optTasks;
  }, [view, optTasks, optArchivedTasks, optFolderTasks]);

  const actualVisibleTasks = useMemo(() => {
    if (searchQuery) {
      let results = searchResults;
      if (searchScope === "current" && view.type === "system") {
        if (view.key === "today") results = results.filter((t) => t.due_at && isToday(new Date(t.due_at)));
        if (view.key === "active") results = results.filter((t) => !t.completed);
        if (view.key === "done") results = results.filter((t) => t.completed);
      }
      return results;
    }
    return visibleTasks;
  }, [searchQuery, searchScope, view, searchResults, visibleTasks]);

  const isArchiveView = view.type === "system" && view.key === "archive";
  const title = searchQuery
    ? "Результаты поиска"
    : view.type === "folder"
      ? view.title
      : SYSTEM_LABELS[view.key];
  const emptyHint = searchQuery
    ? "Ничего не найдено."
    : view.type === "folder"
      ? "В этом списке пока пусто."
      : SYSTEM_EMPTY_HINTS[view.key];

  const selectedTask = useMemo(
    () => actualVisibleTasks.find((t) => t.id === selectedTaskId) ?? null,
    [actualVisibleTasks, selectedTaskId]
  );

  useKeyboardShortcuts({
    tasks: actualVisibleTasks,
    selectedTaskId,
    setSelectedTaskId,
    onArchive: handleArchive,
  });

  return (
    <div className="app-shell">
      <Sidebar
        view={view}
        onSelectSystem={handleSelectSystem}
        counts={counts}
        folders={folders}
        foldersLoading={foldersLoading}
        onSelectFolder={handleSelectFolder}
        onCreateFolder={handleCreateFolder}
        onDeleteFolder={handleDeleteFolder}
        onMoveToFolder={handleMoveToFolder}
        user={user}
        theme={theme}
        onToggleTheme={toggleTheme}
        onLogout={logout}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        searchScope={searchScope}
        setSearchScope={setSearchScope}
        showScopeToggle={view.type === "folder" || view.type === "system"}
      />

      <main className="app-main">
        <TaskList
          title={title}
          emptyHint={emptyHint}
          tasks={actualVisibleTasks}
          loading={searchQuery ? isSearching : loading}
          isArchiveView={isArchiveView}
          canAdd={!isArchiveView && !searchQuery}
          contextLabel={view.type === "folder" ? view.title : null}
          allTags={tags}
          selectedTaskId={selectedTaskId}
          batchSelectedIds={batchSelectedIds}
          onToggleBatchSelect={handleToggleBatchSelect}
          onBulkComplete={handleBulkComplete}
          onBulkArchive={handleBulkArchive}
          onBulkMove={handleBulkMove}
          folders={folders}
          onAdd={handleAdd}
          onToggle={handleToggle}
          onSelect={(task) => setSelectedTaskId(task.id)}
          onArchive={handleArchive}
          onReorder={searchQuery ? null : handleReorderTask}
        />
      </main>

      {selectedTask && (
        <TaskDetail
          key={selectedTask.id}
          task={selectedTask}
          allTags={tags}
          onClose={() => setSelectedTaskId(null)}
          onPatch={handleEdit}
          onArchive={handleArchive}
          onUnarchive={handleUnarchive}
          onPermanentlyDelete={handlePermanentlyDelete}
        />
      )}

      <Toast message={error} onDismiss={() => setError(null)} />
    </div>
  );
}

function Root() {
  const { status } = useAuth();

  if (status === "loading") {
    return (
      <div className="app-shell app-shell--centered">
        <span className="boot-line">Загрузка…</span>
      </div>
    );
  }

  if (status === "guest") {
    return <AuthScreen />;
  }

  return <TodoApp />;
}

export default function App() {
  return (
    <AuthProvider>
      <Root />
    </AuthProvider>
  );
}
