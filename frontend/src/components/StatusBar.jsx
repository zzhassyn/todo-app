export default function StatusBar({ counts, filter, onFilterChange, user, onLogout }) {
  return (
    <footer className="statusbar">
      <div className="statusbar__group">
        <span className="statusbar__mode">NORMAL</span>
        <span className="statusbar__sep">·</span>
        <span className="statusbar__path">~/todo/buffer</span>
      </div>

      <div className="statusbar__filters" role="group" aria-label="фильтр задач">
        {[
          ["all", "все", counts.all],
          ["active", "активные", counts.active],
          ["done", "выполнено", counts.done],
        ].map(([key, label, count]) => (
          <button
            key={key}
            type="button"
            className={`statusbar__filter ${filter === key ? "is-active" : ""}`}
            onClick={() => onFilterChange(key)}
          >
            {label} <span className="statusbar__count">{count}</span>
          </button>
        ))}
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
