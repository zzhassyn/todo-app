ALTER TABLE todoapp.tasks
    ADD COLUMN priority VARCHAR(10) NOT NULL DEFAULT 'medium'
        CHECK (priority IN ('low', 'medium', 'high')),
    ADD COLUMN due_at TIMESTAMPTZ,
    ADD COLUMN archived_at TIMESTAMPTZ;

CREATE TABLE todoapp.tags (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(50) NOT NULL CHECK (char_length(name) BETWEEN 1 AND 50),
    UNIQUE (name)
);

CREATE TABLE todoapp.task_tags (
    task_id INTEGER NOT NULL REFERENCES todoapp.tasks(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES todoapp.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, tag_id)
);

CREATE INDEX task_tags_tag_id_idx ON todoapp.task_tags (tag_id);
