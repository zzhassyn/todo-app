import { useCallback, useEffect, useMemo, useState } from "react";
import { AuthProvider, useAuth } from "./context/AuthContext";
import { api, ApiError } from "./api/client";
import { isToday } from "./utils/date";
import { useTheme } from "./utils/useTheme";
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

  const [tasks, setTasks] = useState([]);
  const [archivedTasks, setArchivedTasks] = useState([]);
  const [archiveLoaded, setArchiveLoaded] = useState(false);
  const [folderTasks, setFolderTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [searchQuery, setSearchQuery] = useState("");
  const [searchScope, setSearchScope] = useState("all");
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);

  const [folders, setFolders] = useState([]);
  const [foldersLoading, setFoldersLoading] = useState(true);
  const [tags, setTags] = useState([]);

  const reportError = useCallback((err, fallback) => {
    setError(err instanceof ApiError ? err.message : fallback);
  }, []);

  // Unfiled tasks (used by all the system filters except "архив") are
  // fetched once up front; the archive and individual folders are fetched
  // lazily, only when first viewed.
  const fetchTasks = useCallback(async (signal) => {
    setLoading(true);
    try {
      const result = await api.listTasks();
      if (signal?.aborted) return;
      setTasks(result ?? []);
    } catch (err) {
      if (signal?.aborted) return;
      reportError(err, "Не удалось загрузить задачи");
    } finally {
      if (!signal?.aborted) setLoading(false);
    }
  }, [reportError]);

  // fetch-on-mount: there is no external store/cache layer in this small
  // app, so the effect's job IS to kick off the initial load. The abort
  // controller guards against the classic race/unmount issue.
  useEffect(() => {
    const controller = new AbortController();
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchTasks(controller.signal);
    return () => controller.abort();
  }, [fetchTasks]);

  useEffect(() => {
    const controller = new AbortController();

    async function loadFolders() {
      setFoldersLoading(true);
      try {
        const result = await api.listFolders();
        if (controller.signal.aborted) return;
        setFolders(result ?? []);
      } catch (err) {
        if (controller.signal.aborted) return;
        reportError(err, "Не удалось загрузить списки");
      } finally {
        if (!controller.signal.aborted) setFoldersLoading(false);
      }
    }

    loadFolders();
    return () => controller.abort();
  }, [reportError]);

  useEffect(() => {
    const controller = new AbortController();

    async function loadTags() {
      try {
        const result = await api.listTags();
        if (controller.signal.aborted) return;
        setTags(result ?? []);
      } catch {
        // Tags are a convenience (autocomplete) feature; a failed load
        // shouldn't block the rest of the app or show an error toast.
      }
    }

    loadTags();
    return () => controller.abort();
  }, []);

  // The archive tab is fetched lazily on first visit rather than upfront,
  // since most sessions never open it.
  useEffect(() => {
    const isArchiveView = view.type === "system" && view.key === "archive";
    if (!isArchiveView || archiveLoaded) return;
    const controller = new AbortController();

    async function loadArchive() {
      setLoading(true);
      try {
        const result = await api.listTasks({ archived: true });
        if (controller.signal.aborted) return;
        setArchivedTasks(result ?? []);
        setArchiveLoaded(true);
      } catch (err) {
        if (controller.signal.aborted) return;
        reportError(err, "Не удалось загрузить архив");
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }

    loadArchive();
    return () => controller.abort();
  }, [view, archiveLoaded, reportError]);

  useEffect(() => {
    if (!searchQuery) {
      setSearchResults([]);
      return;
    }
    const controller = new AbortController();
    
    async function loadSearch() {
      setIsSearching(true);
      try {
        const params = { search: searchQuery };
        if (searchScope === "current") {
          if (view.type === "folder") {
            params.folderId = view.id;
          } else if (view.key === "archive") {
            params.archived = true;
          }
        }
        const result = await api.listTasks(params);
        if (controller.signal.aborted) return;
        setSearchResults(result ?? []);
      } catch (err) {
        if (controller.signal.aborted) return;
        reportError(err, "Ошибка при поиске");
      } finally {
        if (!controller.signal.aborted) setIsSearching(false);
      }
    }
    
    loadSearch();
    return () => controller.abort();
  }, [searchQuery, searchScope, view, reportError]);

  // Each folder's tasks are fetched fresh every time it's selected, rather
  // than cached per-folder: folder contents change often enough (tasks
  // moved in/out) that a per-folder cache would need its own invalidation
  // logic for little benefit at this scale.
  useEffect(() => {
    if (view.type !== "folder") return;
    const controller = new AbortController();

    async function loadFolderTasks() {
      setLoading(true);
      try {
        const result = await api.listTasks({ folderId: view.id });
        if (controller.signal.aborted) return;
        setFolderTasks(result ?? []);
      } catch (err) {
        if (controller.signal.aborted) return;
        reportError(err, "Не удалось загрузить список");
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }

    loadFolderTasks();
    return () => controller.abort();
  }, [view, reportError]);

  const handleSelectSystem = useCallback((key) => {
    setView({ type: "system", key });
    setSelectedTaskId(null);
  }, []);

  const handleSelectFolder = useCallback((folder) => {
    setView({ type: "folder", id: folder.id, title: folder.title });
    setSelectedTaskId(null);
  }, []);

  const handleCreateFolder = useCallback(
    async (title) => {
      try {
        const created = await api.createFolder({ title });
        setFolders((prev) => [...prev, created]);
      } catch (err) {
        reportError(err, "Не удалось создать список");
      }
    },
    [reportError]
  );

  const handleDeleteFolder = useCallback(
    async (id) => {
      try {
        await api.deleteFolder(id);
        setFolders((prev) => prev.filter((f) => f.id !== id));
        // Deleting a folder cascades to its tasks on the backend, so if
        // it's currently open, fall back to the default view rather than
        // showing a now-meaningless empty/stale list.
        setView((prev) => (prev.type === "folder" && prev.id === id ? { type: "system", key: "all" } : prev));
      } catch (err) {
        reportError(err, "Не удалось удалить список");
      }
    },
    [reportError]
  );

  const handleAdd = useCallback(
    async (payload) => {
      try {
        const finalPayload =
          view.type === "folder" ? { ...payload, folder_id: payload.folder_id ?? view.id } : payload;
        const created = await api.createTask(finalPayload);
        if (created.folder_id) {
          setFolderTasks((prev) => [...prev, created]);
        } else {
          setTasks((prev) => [...prev, created]);
        }
      } catch (err) {
        reportError(err, "Не удалось создать задачу");
      }
    },
    [reportError, view]
  );

  const patchInPlace = useCallback((updated) => {
    setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
    setFolderTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
    setArchivedTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
  }, []);

  const handleToggle = useCallback(
    async (task) => {
      try {
        const updated = task.completed
          ? await api.uncompleteTask(task.id)
          : await api.completeTask(task.id);
        patchInPlace(updated);
      } catch (err) {
        reportError(err, "Не удалось изменить статус задачи");
      }
    },
    [reportError, patchInPlace]
  );

  const handleArchive = useCallback(
    async (id) => {
      try {
        const archived = await api.archiveTask(id);
        setTasks((prev) => prev.filter((t) => t.id !== id));
        setFolderTasks((prev) => prev.filter((t) => t.id !== id));
        setArchivedTasks((prev) => (archiveLoaded ? [archived, ...prev] : prev));
        setSelectedTaskId((prev) => (prev === id ? null : prev));
      } catch (err) {
        reportError(err, "Не удалось архивировать задачу");
      }
    },
    [reportError, archiveLoaded]
  );

  const handleUnarchive = useCallback(
    async (id) => {
      try {
        const restored = await api.unarchiveTask(id);
        setArchivedTasks((prev) => prev.filter((t) => t.id !== id));
        if (restored.folder_id) {
          setFolderTasks((prev) => [...prev, restored]);
        } else {
          setTasks((prev) => [...prev, restored]);
        }
      } catch (err) {
        reportError(err, "Не удалось восстановить задачу");
      }
    },
    [reportError]
  );

  const handlePermanentlyDelete = useCallback(
    async (id) => {
      try {
        await api.permanentlyDeleteTask(id);
        setArchivedTasks((prev) => prev.filter((t) => t.id !== id));
        setSelectedTaskId((prev) => (prev === id ? null : prev));
      } catch (err) {
        reportError(err, "Не удалось удалить задачу навсегда");
      }
    },
    [reportError]
  );

  const handleEdit = useCallback(
    async (id, patch) => {
      try {
        const updated = await api.patchTask(id, patch);
        patchInPlace(updated);
      } catch (err) {
        reportError(err, "Не удалось изменить задачу");
      }
    },
    [reportError, patchInPlace]
  );

  const handleMoveToFolder = useCallback(async (taskId, targetFolderId) => {
    const task = tasks.find(t => t.id === taskId) || folderTasks.find(t => t.id === taskId) || archivedTasks.find(t => t.id === taskId);
    if (!task || task.folder_id === targetFolderId) return;

    const backupTasks = [...tasks];
    const backupFolderTasks = [...folderTasks];
    
    const optimisticTask = { ...task, folder_id: targetFolderId };
    
    setTasks(prev => prev.map(t => t.id === task.id ? optimisticTask : t));
    
    if (view.type === "folder" && view.id === task.folder_id) {
      setFolderTasks(prev => prev.filter(t => t.id !== task.id));
    }

    try {
      const updated = await api.patchTask(task.id, { folder_id: targetFolderId });
      patchInPlace(updated);
    } catch (err) {
      reportError(err, "Не удалось переместить задачу");
      setTasks(backupTasks);
      setFolderTasks(backupFolderTasks);
    }
  }, [tasks, folderTasks, archivedTasks, view, patchInPlace, reportError]);

  const handleReorderTask = useCallback(async (draggedId, targetId, insertPosition) => {
    if (draggedId === targetId) return;

    // Use a function to safely get the current visible tasks, since visibleTasks is computed below
    // We can just rely on the dependency array. But to avoid circular deps, let's compute it.
    let currentVisible = tasks;
    if (view.type === "folder") currentVisible = folderTasks;
    else if (view.key === "archive") currentVisible = archivedTasks;
    else if (view.key === "today") currentVisible = tasks.filter((t) => t.due_at && isToday(new Date(t.due_at)));
    else if (view.key === "active") currentVisible = tasks.filter((t) => !t.completed);
    else if (view.key === "done") currentVisible = tasks.filter((t) => t.completed);

    const draggedIdx = currentVisible.findIndex(t => t.id === draggedId);
    const targetIdx = currentVisible.findIndex(t => t.id === targetId);
    if (draggedIdx === -1 || targetIdx === -1) return;

    const draggedTask = currentVisible[draggedIdx];
    
    const newList = [...currentVisible];
    newList.splice(draggedIdx, 1);
    
    const insertIdx = newList.findIndex(t => t.id === targetId);
    const finalIdx = insertPosition === "after" ? insertIdx + 1 : insertIdx;
    
    newList.splice(finalIdx, 0, draggedTask);

    const prevTask = newList[finalIdx - 1];
    const nextTask = newList[finalIdx + 1];

    let newPosition = 0;
    if (prevTask && nextTask) {
      newPosition = (prevTask.position + nextTask.position) / 2.0;
    } else if (prevTask) {
      newPosition = prevTask.position + 1024.0;
    } else if (nextTask) {
      newPosition = nextTask.position - 1024.0;
    } else {
      newPosition = Date.now() / 1000.0;
    }

    const optimisticTask = { ...draggedTask, position: newPosition };
    // Optimistically update all lists to sort it
    const sortByPosition = (a, b) => {
      // also keep folder grouping if needed, but in views they are already filtered by folder
      return a.position - b.position || a.id - b.id;
    };

    setTasks(prev => prev.map(t => t.id === draggedId ? optimisticTask : t).sort(sortByPosition));
    setFolderTasks(prev => prev.map(t => t.id === draggedId ? optimisticTask : t).sort(sortByPosition));
    setArchivedTasks(prev => prev.map(t => t.id === draggedId ? optimisticTask : t).sort(sortByPosition));

    try {
      const updated = await api.patchTask(draggedId, { position: newPosition });
      patchInPlace(updated);
    } catch (err) {
      reportError(err, "Не удалось изменить порядок задач");
      patchInPlace(draggedTask);
    }
  }, [view, tasks, folderTasks, archivedTasks, patchInPlace, reportError]);

  const counts = useMemo(() => {
    const done = tasks.filter((t) => t.completed).length;
    const today = tasks.filter((t) => t.due_at && isToday(new Date(t.due_at))).length;
    return {
      all: tasks.length,
      today,
      done,
      active: tasks.length - done,
      archive: archivedTasks.length,
    };
  }, [tasks, archivedTasks]);

  const visibleTasks = useMemo(() => {
    if (view.type === "folder") return folderTasks;
    if (view.key === "archive") return archivedTasks;
    if (view.key === "today") return tasks.filter((t) => t.due_at && isToday(new Date(t.due_at)));
    if (view.key === "active") return tasks.filter((t) => !t.completed);
    if (view.key === "done") return tasks.filter((t) => t.completed);
    return tasks;
  }, [view, tasks, archivedTasks, folderTasks]);

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
  const title = searchQuery ? "Результаты поиска" : (view.type === "folder" ? view.title : SYSTEM_LABELS[view.key]);
  const emptyHint = searchQuery ? "Ничего не найдено." : (view.type === "folder" ? "В этом списке пока пусто." : SYSTEM_EMPTY_HINTS[view.key]);

  const selectedTask = useMemo(
    () => actualVisibleTasks.find((t) => t.id === selectedTaskId) ?? null,
    [actualVisibleTasks, selectedTaskId]
  );

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
      />

      <main className="app-main">
        <SearchBar
          value={searchQuery}
          onChange={setSearchQuery}
          scope={searchScope}
          onScopeChange={setSearchScope}
          showScopeToggle={view.type === "folder" || view.type === "system"}
        />
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
