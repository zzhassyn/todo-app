import { useCallback, useEffect, useMemo, useState } from "react";
import { AuthProvider, useAuth } from "./context/AuthContext";
import { api, ApiError } from "./api/client";
import AuthScreen from "./components/AuthScreen";
import Buffer from "./components/Buffer";
import StatusBar from "./components/StatusBar";
import Toast from "./components/Toast";
import "./app.css";

function TodoApp() {
  const { user, logout } = useAuth();
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("all");
  const [error, setError] = useState(null);

  const reportError = useCallback((err, fallback) => {
    setError(err instanceof ApiError ? err.message : fallback);
  }, []);

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

  const handleAdd = useCallback(
    async (title) => {
      try {
        const created = await api.createTask({ title });
        setTasks((prev) => [...prev, created]);
      } catch (err) {
        reportError(err, "не удалось создать задачу");
      }
    },
    [reportError]
  );

  const handleToggle = useCallback(
    async (task) => {
      try {
        const updated = task.completed
          ? await api.uncompleteTask(task.id)
          : await api.completeTask(task.id);
        setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
      } catch (err) {
        reportError(err, "не удалось изменить статус задачи");
      }
    },
    [reportError]
  );

  const handleDelete = useCallback(
    async (id) => {
      try {
        await api.deleteTask(id);
        setTasks((prev) => prev.filter((t) => t.id !== id));
      } catch (err) {
        reportError(err, "не удалось удалить задачу");
      }
    },
    [reportError]
  );

  const handleEdit = useCallback(
    async (id, title) => {
      try {
        const updated = await api.patchTask(id, { title });
        setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
      } catch (err) {
        reportError(err, "не удалось изменить задачу");
      }
    },
    [reportError]
  );

  const counts = useMemo(() => {
    const done = tasks.filter((t) => t.completed).length;
    return { all: tasks.length, done, active: tasks.length - done };
  }, [tasks]);

  const visibleTasks = useMemo(() => {
    if (filter === "active") return tasks.filter((t) => !t.completed);
    if (filter === "done") return tasks.filter((t) => t.completed);
    return tasks;
  }, [tasks, filter]);

  return (
    <div className="app-shell">
      <main className="app-main">
        <Buffer
          tasks={visibleTasks}
          loading={loading}
          filter={filter}
          onAdd={handleAdd}
          onToggle={handleToggle}
          onDelete={handleDelete}
          onEdit={handleEdit}
        />
      </main>

      <StatusBar
        counts={counts}
        filter={filter}
        onFilterChange={setFilter}
        user={user}
        onLogout={logout}
      />

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
