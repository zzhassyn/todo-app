import React, { useState, useEffect } from "react";

export default function SearchBar({ value, onChange, scope, onScopeChange, showScopeToggle }) {
  const [localValue, setLocalValue] = useState(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      onChange(localValue);
    }, 300);

    return () => clearTimeout(handler);
  }, [localValue, onChange]);

  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  return (
    <div className="search-bar">
      <div className="search-bar__input-wrapper">
        <svg
          className="search-bar__icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
        <input
          type="text"
          className="search-bar__input"
          placeholder="Поиск задач..."
          value={localValue}
          onChange={(e) => setLocalValue(e.target.value)}
        />
        {localValue && (
          <button
            className="search-bar__clear"
            onClick={() => setLocalValue("")}
            title="Очистить"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        )}
      </div>

      {showScopeToggle && localValue && (
        <div className="search-bar__scope">
          <label className="search-bar__scope-label">
            <input
              type="radio"
              name="searchScope"
              value="current"
              checked={scope === "current"}
              onChange={() => onScopeChange("current")}
            />
            В текущем списке
          </label>
          <label className="search-bar__scope-label">
            <input
              type="radio"
              name="searchScope"
              value="all"
              checked={scope === "all"}
              onChange={() => onScopeChange("all")}
            />
            Везде
          </label>
        </div>
      )}
    </div>
  );
}
