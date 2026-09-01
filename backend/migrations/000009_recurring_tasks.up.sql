CREATE TABLE todoapp.recurring_tasks (
    id SERIAL PRIMARY KEY,
    author_user_id INTEGER NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    priority VARCHAR(10) NOT NULL,
    folder_id UUID REFERENCES todoapp.folders(id) ON DELETE SET NULL,
    tags TEXT[],
    cron_expression VARCHAR(50) NOT NULL,
    last_run_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX recurring_tasks_next_run_at_idx ON todoapp.recurring_tasks (next_run_at);
