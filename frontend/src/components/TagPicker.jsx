import { useEffect, useMemo, useRef, useState } from "react";

/**
 * `allTags` is the full tag vocabulary (from GET /tags), `selected` is the
 * list of tag names already attached to the task being edited. Selecting
 * an existing tag or typing a new name and pressing Enter both call
 * onToggle with a plain string name — the caller doesn't need to
 * distinguish "existing" from "new" since the backend already does
 * find-or-create.
 */
export default function TagPicker({ allTags, selected, onToggle, onClose, anchorRef }) {
  const [query, setQuery] = useState("");
  const popoverRef = useRef(null);
  const inputRef = useRef(null);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  useEffect(() => {
    function handleClickOutside(e) {
      if (popoverRef.current?.contains(e.target)) return;
      if (anchorRef?.current?.contains(e.target)) return;
      onClose();
    }
    function handleKey(e) {
      if (e.key === "Escape") onClose();
    }
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleKey);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKey);
    };
  }, [onClose, anchorRef]);

  const normalizedQuery = query.trim().toLowerCase();

  const filtered = useMemo(
    () => allTags.filter((t) => t.name.toLowerCase().includes(normalizedQuery)),
    [allTags, normalizedQuery]
  );

  const exactMatchExists = allTags.some((t) => t.name.toLowerCase() === normalizedQuery);
  const canCreate = normalizedQuery.length > 0 && !exactMatchExists;

  function handleKeyDown(e) {
    if (e.key === "Enter" && canCreate) {
      e.preventDefault();
      onToggle(normalizedQuery);
      setQuery("");
    }
  }

  return (
    <div className="tagpicker" ref={popoverRef} role="dialog" aria-label="Выбор тегов">
      <input
        ref={inputRef}
        className="tagpicker__input"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="Найти или создать тег…"
        maxLength={50}
      />

      <div className="tagpicker__list">
        {filtered.length === 0 && !canCreate && (
          <div className="tagpicker__empty">Тегов пока нет</div>
        )}

        {filtered.map((tag) => {
          const isSelected = selected.includes(tag.name);
          return (
            <button
              key={tag.id}
              type="button"
              className={`tagpicker__option ${isSelected ? "is-selected" : ""}`}
              onClick={() => onToggle(tag.name)}
            >
              <span className="tagpicker__checkbox" aria-hidden="true">
                {isSelected ? "✓" : ""}
              </span>
              {tag.name}
            </button>
          );
        })}

        {canCreate && (
          <button
            type="button"
            className="tagpicker__option tagpicker__option--create"
            onClick={() => {
              onToggle(normalizedQuery);
              setQuery("");
            }}
          >
            Создать «{normalizedQuery}»
          </button>
        )}
      </div>
    </div>
  );
}
