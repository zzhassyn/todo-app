import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";

export function useTaskSearch({ view, reportError }) {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchScope, setSearchScope] = useState("all");
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    if (!searchQuery) {
      setSearchResults([]);
      return;
    }
    const controller = new AbortController();

    async function loadSearch() {
      setIsSearching(true);
      try {
        let results = [];
        if (searchScope === "all") {
          const [active, archived] = await Promise.all([
            api.listTasks({ search: searchQuery }),
            api.listTasks({ search: searchQuery, archived: true }),
          ]);
          results = [...(active ?? []), ...(archived ?? [])];
        } else {
          const params = { search: searchQuery };
          if (view.type === "folder") {
            params.folderId = view.id;
          } else if (view.key === "archive") {
            params.archived = true;
          }
          results = (await api.listTasks(params)) ?? [];
        }
        
        if (controller.signal.aborted) return;
        setSearchResults(results);
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

  return {
    searchQuery,
    setSearchQuery,
    searchScope,
    setSearchScope,
    searchResults,
    isSearching,
  };
}
