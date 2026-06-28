export default function StatusBar({ mode, user, onLogout }) {
  return (
    <footer className="statusbar">
      <div className="statusbar__group">
        <span className={`statusbar__mode ${mode === "archive" ? "statusbar__mode--archive" : ""}`}>
          {mode === "archive" ? "ARCHIVE" : "NORMAL"}
        </span>
      </div>

      <div className="statusbar__group statusbar__group--right">
        <span className="statusbar__user">{user?.email}</span>
        <button type="button" className="statusbar__logout" onClick={onLogout} title="выйти">
          :q
        </button>
      </div>
    </footer>
  );
}
