import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";

export function useFolders({ reportError }) {
  const [folders, setFolders] = useState([]);
  const [foldersLoading, setFoldersLoading] = useState(true);

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
      } catch (err) {
        reportError(err, "Не удалось удалить список");
      }
    },
    [reportError]
  );

  return {
    folders,
    foldersLoading,
    handleCreateFolder,
    handleDeleteFolder,
  };
}
