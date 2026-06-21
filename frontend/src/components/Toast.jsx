export default function Toast({ message, onDismiss }) {
  if (!message) return null;

  return (
    <div className="toast" role="alert">
      <span className="toast__badge">error</span>
      <span className="toast__message">{message}</span>
      <button type="button" className="toast__dismiss" onClick={onDismiss} aria-label="закрыть">
        ×
      </button>
    </div>
  );
}
