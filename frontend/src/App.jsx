import { useCallback, useEffect, useMemo, useState } from "react";
import { AuthProvider, useAuth } from "./context/AuthContext";
import { api, ApiError } from "./api/client";
import AuthScreen from "./components/AuthScreen";
import Sidebar from "./components/Sidebar";
import Buffer from "./components/Buffer";
import StatusBar from "./components/StatusBar";
import Toast from "./components/Toast";
import "./app.css";

const SYSTEM_LABELS = {
  all: "buffer.todo",
  active: "active.todo",
  done: "done.todo",
  archive: "archive.todo",
};

const SYSTEM_EMPTY_HINTS = {
  all: "буфер пуст. начните печатать ниже.",
  active: "активных задач нет — всё сделано.",
  done: "ничего не выполнено. пока.",
  archive: "архив пуст.",
};

function TodoApp() {
  const { user, logout } = useAuth();

  // `view` is the single source of truth for "what is currently shown":
  // either one of the system filters, or a specific folder. Keeping it as
  // one discriminated value (rather than separate `filter` + `folderId`
  // state) means there's exactly one place that decides what's selected,
  // so the sidebar's highlighted item and the buffer's contents can never
  // disagree with each other.
  const [view, setView] = useState({ type: "system", key: "all" });

  const [tasks, setTasks] = useState([]);
  const [archivedTasks, setArchivedTasks] = useState([]);
  const [archiveLoaded, setArchiveLoaded] = useState(false);
  const [folderTasks, setFolderTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [folders, setFolders] = useState([]);
  const [foldersLoading, setFoldersLoading] = useState(true);

  const reportError = useCallback((err, fallback) => {
    setError(err instanceof ApiError ? err.message : fallback);
  }, []);

  // Unfiled tasks (used by the "все"/"активные"/"выполнено" system
  // filters) are fetched once up front; the archive and individual
  // folders are fetched lazily, only when first viewed.
  const fetchTasks = useCallback(async (signal) => {
    setLoading(true);
    try {
      const result = await api.listTasks();
      if (signal?.aborted) return;
      setTasks(result ?? []);
    } catch (err) {
      if (signal?.aborted) return;
      reportError(err, "не удалось загрузить задачи");
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
        reportError(err, "не удалось загрузить списки");
      } finally {
        if (!controller.signal.aborted) setFoldersLoading(false);
      }
    }

    loadFolders();
    return () => controller.abort();
  }, [reportError]);

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
        reportError(err, "не удалось загрузить архив");
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }

    loadArchive();
    return () => controller.abort();
  }, [view, archiveLoaded, reportError]);

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
        reportError(err, "не удалось загрузить список");
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    }

    loadFolderTasks();
    return () => controller.abort();
  }, [view, reportError]);

  const handleSelectSystem = useCallback((key) => {
    setView({ type: "system", key });
  }, []);

  const handleSelectFolder = useCallback((folder) => {
    setView({ type: "folder", id: folder.id, title: folder.title });
  }, []);

  const handleCreateFolder = useCallback(
    async (title) => {
      try {
        const created = await api.createFolder({ title });
        setFolders((prev) => [...prev, created]);
      } catch (err) {
        reportError(err, "не удалось создать список");
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
        reportError(err, "не удалось удалить список");
      }
    },
    [reportError]
  );

  const handleAdd = useCallback(
    async (payload) => {
      try {
        const created = await api.createTask(payload);
        if (created.folder_id) {
          setFolderTasks((prev) => [...prev, created]);
        } else {
          setTasks((prev) => [...prev, created]);
        }
      } catch (err) {
        reportError(err, "не удалось создать задачу");
      }
    },
    [reportError]
  );

  const patchInPlace = useCallback((updated) => {
    setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
    setFolderTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
  }, []);

  const handleToggle = useCallback(
    async (task) => {
      try {
        const updated = task.completed
          ? await api.uncompleteTask(task.id)
          : await api.completeTask(task.id);
        patchInPlace(updated);
      } catch (err) {
        reportError(err, "не удалось изменить статус задачи");
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
      } catch (err) {
        reportError(err, "не удалось архивировать задачу");
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
        reportError(err, "не удалось восстановить задачу");
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
        reportError(err, "не удалось изменить задачу");
      }
    },
    [reportError, patchInPlace]
  );

  const counts = useMemo(() => {
    const done = tasks.filter((t) => t.completed).length;
    return { all: tasks.length, done, active: tasks.length - done, archive: archivedTasks.length };
  }, [tasks, archivedTasks]);

  const visibleTasks = useMemo(() => {
    if (view.type === "folder") return folderTasks;
    if (view.key === "archive") return archivedTasks;
    if (view.key === "active") return tasks.filter((t) => !t.completed);
    if (view.key === "done") return tasks.filter((t) => t.completed);
    return tasks;
  }, [view, tasks, archivedTasks, folderTasks]);

  const title = view.type === "folder" ? `${view.title}.todo` : SYSTEM_LABELS[view.key];
  const emptyHint = view.type === "folder" ? "в этом списке пока пусто." : SYSTEM_EMPTY_HINTS[view.key];

  return (
    <div className="app-shell app-shell--with-sidebar">
      <Sidebar
        view={view}
        onSelectSystem={handleSelectSystem}
        counts={counts}
        folders={folders}
        foldersLoading={foldersLoading}
        onSelectFolder={handleSelectFolder}
        onCreateFolder={handleCreateFolder}
        onDeleteFolder={handleDeleteFolder}
      />

      <div className="app-shell__content">
        <main className="app-main">
          <Buffer
            view={view}
            title={title}
            emptyHint={emptyHint}
            tasks={visibleTasks}
            loading={loading}
            onAdd={handleAdd}
            onToggle={handleToggle}
            onArchive={handleArchive}
            onUnarchive={handleUnarchive}
            onEdit={handleEdit}
          />
        </main>

        <StatusBar
          mode={view.type === "system" && view.key === "archive" ? "archive" : "normal"}
          user={user}
          onLogout={logout}
        />
      </div>

      <Toast message={error} onDismiss={() => setError(null)} />
    </div>
  );
}

function Root() {
  const { status } = useAuth();

  if (status === "loading") {
    return (
      <div className="app-shell app-shell--centered">
        <span className="boot-line">инициализация буфера…</span>
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
