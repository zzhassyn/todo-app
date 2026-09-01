import { useCallback, useEffect, useState, useOptimistic, startTransition } from "react";
import { api, ApiError } from "../api/client";
import { isToday } from "../utils/date";

export function useTasks({ view, reportError }) {
  const [tasks, setTasks] = useState([]);
  const [archivedTasks, setArchivedTasks] = useState([]);
  const [archiveLoaded, setArchiveLoaded] = useState(false);
  const [folderTasks, setFolderTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [tags, setTags] = useState([]);

  const [batchSelectedIds, setBatchSelectedIds] = useState(new Set());

  const taskReducer = useCallback((state, action) => {
    if (action.type === "patch") {
      return state.map((t) => (t.id === action.payload.id ? { ...t, ...action.payload } : t));
    }
    if (action.type === "remove") {
      const idSet = new Set(action.payload);
      return state.filter((t) => !idSet.has(t.id));
    }
    if (action.type === "add") {
      return [...state, ...action.payload];
    }
    return state;
  }, []);

  const [optTasks, setOptTasks] = useOptimistic(tasks, taskReducer);
  const [optFolderTasks, setOptFolderTasks] = useOptimistic(folderTasks, taskReducer);
  const [optArchivedTasks, setOptArchivedTasks] = useOptimistic(archivedTasks, taskReducer);

  const fetchTasks = useCallback(
    async (signal) => {
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
    },
    [reportError]
  );

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

  const patchInPlace = useCallback((updated) => {
    setTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
    setFolderTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
    setArchivedTasks((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
  }, []);

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

  const handleToggle = useCallback(
    async (task) => {
      startTransition(() => {
        setOptTasks({ type: "patch", payload: { id: task.id, completed: !task.completed } });
        setOptFolderTasks({ type: "patch", payload: { id: task.id, completed: !task.completed } });
        setOptArchivedTasks({ type: "patch", payload: { id: task.id, completed: !task.completed } });
      });
      try {
        const updated = task.completed ? await api.uncompleteTask(task.id) : await api.completeTask(task.id);
        patchInPlace(updated);
      } catch (err) {
        reportError(err, "Не удалось изменить статус задачи");
      }
    },
    [reportError, patchInPlace, setOptTasks, setOptFolderTasks, setOptArchivedTasks]
  );

  const handleArchive = useCallback(
    async (id) => {
      startTransition(() => {
        setOptTasks({ type: "remove", payload: [id] });
        setOptFolderTasks({ type: "remove", payload: [id] });
      });
      try {
        const archived = await api.archiveTask(id);
        setTasks((prev) => prev.filter((t) => t.id !== id));
        setFolderTasks((prev) => prev.filter((t) => t.id !== id));
        setArchivedTasks((prev) => (archiveLoaded ? [archived, ...prev] : prev));
      } catch (err) {
        reportError(err, "Не удалось архивировать задачу");
      }
    },
    [reportError, archiveLoaded, setOptTasks, setOptFolderTasks]
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

  const handleMoveToFolder = useCallback(
    async (taskId, targetFolderId) => {
      const task =
        tasks.find((t) => t.id === taskId) ||
        folderTasks.find((t) => t.id === taskId) ||
        archivedTasks.find((t) => t.id === taskId);
      if (!task || task.folder_id === targetFolderId) return;

      const backupTasks = [...tasks];
      const backupFolderTasks = [...folderTasks];

      const optimisticTask = { ...task, folder_id: targetFolderId };

      setTasks((prev) => prev.map((t) => (t.id === task.id ? optimisticTask : t)));

      if (view.type === "folder" && view.id === task.folder_id) {
        setFolderTasks((prev) => prev.filter((t) => t.id !== task.id));
      }

      try {
        const updated = await api.patchTask(task.id, { folder_id: targetFolderId });
        patchInPlace(updated);
      } catch (err) {
        reportError(err, "Не удалось переместить задачу");
        setTasks(backupTasks);
        setFolderTasks(backupFolderTasks);
      }
    },
    [tasks, folderTasks, archivedTasks, view, patchInPlace, reportError]
  );

  const handleReorderTask = useCallback(
    async (draggedId, targetId, insertPosition) => {
      if (draggedId === targetId) return;

      let currentVisible = tasks;
      if (view.type === "folder") currentVisible = folderTasks;
      else if (view.key === "archive") currentVisible = archivedTasks;
      else if (view.key === "today")
        currentVisible = tasks.filter((t) => t.due_at && isToday(new Date(t.due_at)));
      else if (view.key === "active") currentVisible = tasks.filter((t) => !t.completed);
      else if (view.key === "done") currentVisible = tasks.filter((t) => t.completed);

      const draggedIdx = currentVisible.findIndex((t) => t.id === draggedId);
      const targetIdx = currentVisible.findIndex((t) => t.id === targetId);
      if (draggedIdx === -1 || targetIdx === -1) return;

      const draggedTask = currentVisible[draggedIdx];

      const newList = [...currentVisible];
      newList.splice(draggedIdx, 1);

      const insertIdx = newList.findIndex((t) => t.id === targetId);
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
      const sortByPosition = (a, b) => a.position - b.position || a.id - b.id;

      setTasks((prev) =>
        prev.map((t) => (t.id === draggedId ? optimisticTask : t)).sort(sortByPosition)
      );
      setFolderTasks((prev) =>
        prev.map((t) => (t.id === draggedId ? optimisticTask : t)).sort(sortByPosition)
      );
      setArchivedTasks((prev) =>
        prev.map((t) => (t.id === draggedId ? optimisticTask : t)).sort(sortByPosition)
      );

      try {
        const updated = await api.patchTask(draggedId, { position: newPosition });
        patchInPlace(updated);
      } catch (err) {
        reportError(err, "Не удалось изменить порядок задач");
        patchInPlace(draggedTask);
      }
    },
    [view, tasks, folderTasks, archivedTasks, patchInPlace, reportError]
  );

  const handleBulkComplete = useCallback(async () => {
    const ids = Array.from(batchSelectedIds);
    if (ids.length === 0) return;
    try {
      const updatedList = await api.bulkCompleteTasks(ids);
      for (const t of updatedList?.tasks || []) patchInPlace(t);
      setBatchSelectedIds(new Set());
    } catch (err) {
      reportError(err, "Не удалось выполнить задачи");
    }
  }, [batchSelectedIds, patchInPlace, reportError]);

  const handleBulkArchive = useCallback(async () => {
    const ids = Array.from(batchSelectedIds);
    if (ids.length === 0) return;
    try {
      const updatedList = await api.bulkArchiveTasks(ids);
      for (const t of updatedList?.tasks || []) patchInPlace(t);

      const idSet = new Set(ids);
      setTasks((prev) => prev.filter((t) => !idSet.has(t.id)));
      setFolderTasks((prev) => prev.filter((t) => !idSet.has(t.id)));

      if (archiveLoaded) {
        setArchivedTasks((prev) => [
          ...(updatedList?.tasks || []),
          ...prev.filter((t) => !idSet.has(t.id)),
        ]);
      }

      setBatchSelectedIds(new Set());
    } catch (err) {
      reportError(err, "Не удалось архивировать задачи");
    }
  }, [batchSelectedIds, patchInPlace, reportError, archiveLoaded]);

  const handleBulkMove = useCallback(
    async (folderId) => {
      const ids = Array.from(batchSelectedIds);
      if (ids.length === 0) return;
      try {
        const updatedList = await api.bulkPatchTasks({ task_ids: ids, folder_id: folderId });
        for (const t of updatedList?.tasks || []) patchInPlace(t);

        const idSet = new Set(ids);
        if (view.type === "folder" && view.id !== folderId) {
          setFolderTasks((prev) => prev.filter((t) => !idSet.has(t.id)));
        }
        setBatchSelectedIds(new Set());
      } catch (err) {
        reportError(err, "Не удалось переместить задачи");
      }
    },
    [batchSelectedIds, patchInPlace, reportError, view]
  );

  const handleToggleBatchSelect = useCallback((id) => {
    setBatchSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }, []);

  return {
    tasks,
    folderTasks,
    archivedTasks,
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
  };
}
