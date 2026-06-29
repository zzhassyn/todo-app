export default function PriorityDot({ priority }) {
  if (!priority) return null;
  return <span className={`priority-dot priority-dot--${priority}`} aria-hidden="true" />;
}
