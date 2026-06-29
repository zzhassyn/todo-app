const WEEKDAY_LABELS = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];
const MONTH_LABELS = [
  "Январь",
  "Февраль",
  "Март",
  "Апрель",
  "Май",
  "Июнь",
  "Июль",
  "Август",
  "Сентябрь",
  "Октябрь",
  "Ноябрь",
  "Декабрь",
];

export function weekdayLabels() {
  return WEEKDAY_LABELS;
}

export function monthLabel(date) {
  return `${MONTH_LABELS[date.getMonth()]} ${date.getFullYear()}`;
}

export function startOfDay(date) {
  const d = new Date(date);
  d.setHours(0, 0, 0, 0);
  return d;
}

export function isSameDay(a, b) {
  if (!a || !b) return false;
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

export function isToday(date) {
  return isSameDay(date, new Date());
}

export function isOverdue(date) {
  return startOfDay(date) < startOfDay(new Date());
}

/**
 * Builds a 6-row x 7-col grid of Dates covering the given month, starting
 * on Monday, including the leading/trailing days from neighboring months
 * needed to fill the grid (standard calendar-widget behavior).
 */
export function buildMonthGrid(year, month) {
  const firstOfMonth = new Date(year, month, 1);
  // getDay(): 0=Sun..6=Sat. We want Monday-first, so shift.
  const leadingOffset = (firstOfMonth.getDay() + 6) % 7;

  const gridStart = new Date(year, month, 1 - leadingOffset);

  const days = [];
  for (let i = 0; i < 42; i++) {
    const d = new Date(gridStart);
    d.setDate(gridStart.getDate() + i);
    days.push(d);
  }
  return days;
}

/** Short, human label: "Сегодня", "Завтра", or "28 июн". Used in compact UI like task rows. */
const SHORT_MONTHS = [
  "янв",
  "фев",
  "мар",
  "апр",
  "май",
  "июн",
  "июл",
  "авг",
  "сен",
  "окт",
  "ноя",
  "дек",
];

export function formatShortDate(date) {
  const today = startOfDay(new Date());
  const target = startOfDay(date);
  const diffDays = Math.round((target - today) / 86400000);

  if (diffDays === 0) return "Сегодня";
  if (diffDays === 1) return "Завтра";
  if (diffDays === -1) return "Вчера";

  return `${date.getDate()} ${SHORT_MONTHS[date.getMonth()]}`;
}

export function toDateInputValue(date) {
  const d = startOfDay(date);
  const yyyy = d.getFullYear();
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const dd = String(d.getDate()).padStart(2, "0");
  return `${yyyy}-${mm}-${dd}`;
}
