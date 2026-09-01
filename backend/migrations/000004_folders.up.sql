CREATE TABLE todoapp.folders (
    id          UUID PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
    title       VARCHAR(100) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX folders_user_id_idx ON todoapp.folders (user_id);

ALTER TABLE todoapp.tasks
    ADD COLUMN folder_id UUID REFERENCES todoapp.folders(id) ON DELETE CASCADE;

CREATE INDEX tasks_folder_id_idx ON todoapp.tasks (folder_id);
