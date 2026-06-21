import { useState } from "react";

export default function Composer({ onAdd, nextIndex, disabled }) {
  const [value, setValue] = useState("");

  function handleSubmit(e) {
    e.preventDefault();
    const trimmed = value.trim();
    if (!trimmed || disabled) return;
    onAdd(trimmed);
    setValue("");
  }

  return (
    <form className="composer" onSubmit={handleSubmit}>
      <span className="line__no line__no--ghost">{String(nextIndex).padStart(2, "0")}</span>
      <span className="line__marker line__marker--add composer__marker">+</span>
      <input
        className="composer__input"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder="новая задача…"
        maxLength={100}
        disabled={disabled}
        aria-label="новая задача"
      />
      <span className="composer__caret" aria-hidden="true" />
    </form>
  );
}
