export default function TagBadge({ name, onRemove }) {
  return (
    <span className="tag-badge">
      {name}
      {onRemove && (
        <button
          type="button"
          className="tag-badge__remove"
          onClick={onRemove}
          aria-label={`Убрать тег ${name}`}
        >
          ×
        </button>
      )}
    </span>
  );
}
