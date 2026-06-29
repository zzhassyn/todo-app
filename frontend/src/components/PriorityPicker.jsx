import { useEffect, useRef } from "react";

const OPTIONS = [
  { value: "low", label: "Низкий" },
  { value: "medium", label: "Средний" },
  { value: "high", label: "Высокий" },
];

export default function PriorityPicker({ value, onChange, onClose, anchorRef }) {
  const popoverRef = useRef(null);

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

  return (
    <div className="prioritypicker" ref={popoverRef} role="listbox" aria-label="Выбор приоритета">
      {OPTIONS.map((opt) => (
        <button
          key={opt.value}
          type="button"
          role="option"
          aria-selected={value === opt.value}
          className={`prioritypicker__option ${value === opt.value ? "is-selected" : ""}`}
          onClick={() => {
            onChange(opt.value);
            onClose();
          }}
        >
          <span className={`priority-dot priority-dot--${opt.value}`} aria-hidden="true" />
          {opt.label}
        </button>
      ))}
    </div>
  );
}
