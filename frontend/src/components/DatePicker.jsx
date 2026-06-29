import { useEffect, useRef, useState } from "react";
import { buildMonthGrid, isSameDay, isToday, monthLabel, weekdayLabels } from "../utils/date";

/**
 * A self-contained calendar popover. Renders inline (not via portal) —
 * the caller is responsible for positioning via a relatively-positioned
 * wrapper, same pattern as the other pickers in this app.
 */
export default function DatePicker({ value, onChange, onClose, anchorRef }) {
  const [cursor, setCursor] = useState(() => value ?? new Date());
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

  const year = cursor.getFullYear();
  const month = cursor.getMonth();
  const days = buildMonthGrid(year, month);

  function goMonth(delta) {
    setCursor(new Date(year, month + delta, 1));
  }

  return (
    <div className="datepicker" ref={popoverRef} role="dialog" aria-label="Выбор даты">
      <div className="datepicker__header">
        <button
          type="button"
          className="datepicker__nav"
          onClick={() => goMonth(-1)}
          aria-label="Предыдущий месяц"
        >
          ‹
        </button>
        <span className="datepicker__label">{monthLabel(cursor)}</span>
        <button
          type="button"
          className="datepicker__nav"
          onClick={() => goMonth(1)}
          aria-label="Следующий месяц"
        >
          ›
        </button>
      </div>

      <div className="datepicker__weekdays">
        {weekdayLabels().map((w) => (
          <span key={w}>{w}</span>
        ))}
      </div>

      <div className="datepicker__grid">
        {days.map((day) => {
          const outside = day.getMonth() !== month;
          const selected = isSameDay(day, value);
          return (
            <button
              key={day.toISOString()}
              type="button"
              className={`datepicker__day ${outside ? "is-outside" : ""} ${
                selected ? "is-selected" : ""
              } ${isToday(day) ? "is-today" : ""}`}
              onClick={() => onChange(day)}
            >
              {day.getDate()}
            </button>
          );
        })}
      </div>

      <div className="datepicker__footer">
        <button
          type="button"
          className="datepicker__quick"
          onClick={() => onChange(new Date())}
        >
          Сегодня
        </button>
        {value && (
          <button
            type="button"
            className="datepicker__quick datepicker__quick--clear"
            onClick={() => onChange(null)}
          >
            Очистить
          </button>
        )}
      </div>
    </div>
  );
}
