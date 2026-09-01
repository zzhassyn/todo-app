import { useEffect } from "react";

/**
 * Registers global keyboard shortcuts for task navigation and actions.
 *
 * - j / k: Navigate down/up through the visible task list
 * - e: Archive the currently selected task
 * - Cmd+K / Ctrl+K: Focus the search bar
 *
 * Shortcuts are suppressed when the user is typing in an input/textarea
 * (except Cmd+K which always works).
 */
export function useKeyboardShortcuts({ tasks, selectedTaskId, setSelectedTaskId, onArchive }) {
  useEffect(() => {
    const handleKeyDown = (e) => {
      // Ignore if typing in an input (except Cmd+K)
      if (
        ["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement.tagName) &&
        !(e.metaKey && e.key === "k")
      ) {
        return;
      }

      if (e.metaKey && e.key === "k") {
        e.preventDefault();
        document.querySelector(".search-bar__input")?.focus();
        return;
      }

      if (e.key === "j" || e.key === "k") {
        e.preventDefault();
        if (tasks.length === 0) return;

        const currentIndex = tasks.findIndex((t) => t.id === selectedTaskId);
        if (e.key === "j") {
          if (currentIndex === -1 || currentIndex === tasks.length - 1) {
            setSelectedTaskId(tasks[0].id);
          } else {
            setSelectedTaskId(tasks[currentIndex + 1].id);
          }
        } else if (e.key === "k") {
          if (currentIndex === -1 || currentIndex === 0) {
            setSelectedTaskId(tasks[tasks.length - 1].id);
          } else {
            setSelectedTaskId(tasks[currentIndex - 1].id);
          }
        }

        // Scroll the selected task into view if needed
        setTimeout(() => {
          document.querySelector(".task-list-item--selected")?.scrollIntoView({
            block: "nearest",
            behavior: "smooth",
          });
        }, 50);
      }

      if (e.key === "e" && selectedTaskId) {
        e.preventDefault();
        onArchive(selectedTaskId);
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [tasks, selectedTaskId, setSelectedTaskId, onArchive]);
}
